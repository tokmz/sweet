package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"sweet/pkg/cache"
)

func main() {
	// 演示单机模式
	demoSingleMode()

	// 演示集群模式 (取消注释以测试)
	// demoClusterMode()

	// 演示哨兵模式 (取消注释以测试)
	// demoSentinelMode()
}

// 演示单机模式
func demoSingleMode() {
	fmt.Println("=== 单机模式示例 ===")

	// 创建配置
	config := cache.Config{
		Mode: cache.ModeSingle,
		Single: struct {
			Addr string `json:"addr"`
		}{
			Addr: "localhost:6379",
		},
		Password:    "",
		DB:          0,
		EnableTrace: true,
	}

	// 创建客户端
	client, err := cache.NewClient(config)
	if err != nil {
		log.Fatalf("创建Redis客户端失败: %v", err)
	}
	defer client.Close()

	// 使用客户端
	ctx := context.Background()

	// 演示字符串操作
	demoStringOperations(ctx, client)

	// 演示哈希表操作
	demoHashOperations(ctx, client)

	// 演示列表操作
	demoListOperations(ctx, client)

	// 演示集合操作
	demoSetOperations(ctx, client)

	// 演示有序集合操作
	demoZSetOperations(ctx, client)

	// 演示自定义命令
	demoCustomCommands(ctx, client)
}

// 演示集群模式
func demoClusterMode() {
	fmt.Println("=== 集群模式示例 ===")

	// 创建配置
	config := cache.Config{
		Mode: cache.ModeCluster,
		Cluster: struct {
			Addrs []string `json:"addrs"`
		}{
			Addrs: []string{
				"localhost:7000",
				"localhost:7001",
				"localhost:7002",
			},
		},
		Password:    "",
		EnableTrace: true,
	}

	// 创建客户端
	client, err := cache.NewClient(config)
	if err != nil {
		log.Fatalf("创建Redis集群客户端失败: %v", err)
	}
	defer client.Close()

	// 使用客户端
	ctx := context.Background()

	// 演示字符串操作
	demoStringOperations(ctx, client)
}

// 演示哨兵模式
func demoSentinelMode() {
	fmt.Println("=== 哨兵模式示例 ===")

	// 创建配置
	config := cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: struct {
			MasterName string   `json:"master_name"`
			Addrs      []string `json:"addrs"`
		}{
			MasterName: "mymaster",
			Addrs: []string{
				"localhost:26379",
				"localhost:26380",
				"localhost:26381",
			},
		},
		Password:    "",
		DB:          0,
		EnableTrace: true,
	}

	// 创建客户端
	client, err := cache.NewClient(config)
	if err != nil {
		log.Fatalf("创建Redis哨兵客户端失败: %v", err)
	}
	defer client.Close()

	// 使用客户端
	ctx := context.Background()

	// 演示字符串操作
	demoStringOperations(ctx, client)
}

// 演示字符串操作
func demoStringOperations(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 字符串操作 ---")

	// 设置键值
	err := client.Set(ctx, "string_key", "Hello Redis", time.Minute)
	if err != nil {
		log.Printf("设置键值失败: %v", err)
		return
	}
	fmt.Println("设置键值成功: string_key = Hello Redis")

	// 获取键值
	val, err := client.Get(ctx, "string_key")
	if err != nil {
		log.Printf("获取键值失败: %v", err)
		return
	}
	fmt.Printf("获取键值成功: string_key = %s\n", val)

	// 设置键值(如果不存在)
	success, err := client.SetNX(ctx, "nx_key", "NX Value", time.Minute)
	if err != nil {
		log.Printf("SetNX失败: %v", err)
		return
	}
	fmt.Printf("SetNX结果: nx_key = %v\n", success)

	// 检查键是否存在
	exists, err := client.Exists(ctx, "string_key")
	if err != nil {
		log.Printf("检查键是否存在失败: %v", err)
		return
	}
	fmt.Printf("键存在检查: string_key exists = %v\n", exists > 0)

	// 设置过期时间
	ok, err := client.Expire(ctx, "string_key", time.Second*30)
	if err != nil {
		log.Printf("设置过期时间失败: %v", err)
		return
	}
	fmt.Printf("设置过期时间: string_key expire = %v\n", ok)

	// 获取剩余生存时间
	ttl, err := client.TTL(ctx, "string_key")
	if err != nil {
		log.Printf("获取TTL失败: %v", err)
		return
	}
	fmt.Printf("剩余生存时间: string_key TTL = %v\n", ttl)

	// 删除键
	deleted, err := client.Del(ctx, "string_key")
	if err != nil {
		log.Printf("删除键失败: %v", err)
		return
	}
	fmt.Printf("删除键: string_key deleted = %v\n", deleted > 0)
}

// 演示哈希表操作
func demoHashOperations(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 哈希表操作 ---")

	// 设置哈希表字段
	n, err := client.HSet(ctx, "hash_key", "field1", "value1", "field2", "value2")
	if err != nil {
		log.Printf("HSet失败: %v", err)
		return
	}
	fmt.Printf("设置哈希表字段: hash_key fields = %d\n", n)

	// 获取哈希表字段值
	val, err := client.HGet(ctx, "hash_key", "field1")
	if err != nil {
		log.Printf("HGet失败: %v", err)
		return
	}
	fmt.Printf("获取哈希表字段值: hash_key[field1] = %s\n", val)

	// 获取所有哈希表字段和值
	all, err := client.HGetAll(ctx, "hash_key")
	if err != nil {
		log.Printf("HGetAll失败: %v", err)
		return
	}
	fmt.Printf("获取所有哈希表字段和值: hash_key = %v\n", all)

	// 检查哈希表字段是否存在
	exists, err := client.HExists(ctx, "hash_key", "field1")
	if err != nil {
		log.Printf("HExists失败: %v", err)
		return
	}
	fmt.Printf("检查哈希表字段是否存在: hash_key[field1] exists = %v\n", exists)

	// 删除哈希表字段
	deleted, err := client.HDel(ctx, "hash_key", "field1")
	if err != nil {
		log.Printf("HDel失败: %v", err)
		return
	}
	fmt.Printf("删除哈希表字段: hash_key[field1] deleted = %v\n", deleted > 0)
}

// 演示列表操作
func demoListOperations(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 列表操作 ---")

	// 从列表左端插入元素
	n, err := client.LPush(ctx, "list_key", "value1", "value2", "value3")
	if err != nil {
		log.Printf("LPush失败: %v", err)
		return
	}
	fmt.Printf("从列表左端插入元素: list_key size = %d\n", n)

	// 从列表右端插入元素
	n, err = client.RPush(ctx, "list_key", "value4", "value5")
	if err != nil {
		log.Printf("RPush失败: %v", err)
		return
	}
	fmt.Printf("从列表右端插入元素: list_key size = %d\n", n)

	// 获取列表长度
	length, err := client.LLen(ctx, "list_key")
	if err != nil {
		log.Printf("LLen失败: %v", err)
		return
	}
	fmt.Printf("获取列表长度: list_key length = %d\n", length)

	// 获取列表元素
	elements, err := client.LRange(ctx, "list_key", 0, -1)
	if err != nil {
		log.Printf("LRange失败: %v", err)
		return
	}
	fmt.Printf("获取列表元素: list_key elements = %v\n", elements)

	// 从左端弹出元素
	val, err := client.LPop(ctx, "list_key")
	if err != nil {
		log.Printf("LPop失败: %v", err)
		return
	}
	fmt.Printf("从左端弹出元素: list_key popped = %s\n", val)

	// 从右端弹出元素
	val, err = client.RPop(ctx, "list_key")
	if err != nil {
		log.Printf("RPop失败: %v", err)
		return
	}
	fmt.Printf("从右端弹出元素: list_key popped = %s\n", val)
}

// 演示集合操作
func demoSetOperations(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 集合操作 ---")

	// 添加集合成员
	n, err := client.SAdd(ctx, "set_key", "member1", "member2", "member3")
	if err != nil {
		log.Printf("SAdd失败: %v", err)
		return
	}
	fmt.Printf("添加集合成员: set_key members = %d\n", n)

	// 获取集合成员数
	count, err := client.SCard(ctx, "set_key")
	if err != nil {
		log.Printf("SCard失败: %v", err)
		return
	}
	fmt.Printf("获取集合成员数: set_key count = %d\n", count)

	// 获取集合所有成员
	members, err := client.SMembers(ctx, "set_key")
	if err != nil {
		log.Printf("SMembers失败: %v", err)
		return
	}
	fmt.Printf("获取集合所有成员: set_key members = %v\n", members)

	// 检查成员是否在集合中
	exists, err := client.SIsMember(ctx, "set_key", "member1")
	if err != nil {
		log.Printf("SIsMember失败: %v", err)
		return
	}
	fmt.Printf("检查成员是否在集合中: set_key has member1 = %v\n", exists)

	// 移除集合成员
	removed, err := client.SRem(ctx, "set_key", "member1")
	if err != nil {
		log.Printf("SRem失败: %v", err)
		return
	}
	fmt.Printf("移除集合成员: set_key removed member1 = %v\n", removed > 0)
}

// 演示有序集合操作
func demoZSetOperations(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 有序集合操作 ---")

	// 添加有序集合成员
	n, err := client.ZAdd(ctx, "zset_key",
		&cache.Z{Score: 1.0, Member: "member1"},
		&cache.Z{Score: 2.0, Member: "member2"},
		&cache.Z{Score: 3.0, Member: "member3"},
	)
	if err != nil {
		log.Printf("ZAdd失败: %v", err)
		return
	}
	fmt.Printf("添加有序集合成员: zset_key members = %d\n", n)

	// 获取有序集合成员数
	count, err := client.ZCard(ctx, "zset_key")
	if err != nil {
		log.Printf("ZCard失败: %v", err)
		return
	}
	fmt.Printf("获取有序集合成员数: zset_key count = %d\n", count)

	// 获取有序集合成员
	members, err := client.ZRange(ctx, "zset_key", 0, -1)
	if err != nil {
		log.Printf("ZRange失败: %v", err)
		return
	}
	fmt.Printf("获取有序集合成员: zset_key members = %v\n", members)

	// 获取有序集合成员及分数
	membersWithScores, err := client.ZRangeWithScores(ctx, "zset_key", 0, -1)
	if err != nil {
		log.Printf("ZRangeWithScores失败: %v", err)
		return
	}
	fmt.Println("获取有序集合成员及分数:")
	for _, z := range membersWithScores {
		fmt.Printf("  zset_key member = %v, score = %f\n", z.Member, z.Score)
	}

	// 移除有序集合成员
	removed, err := client.ZRem(ctx, "zset_key", "member1")
	if err != nil {
		log.Printf("ZRem失败: %v", err)
		return
	}
	fmt.Printf("移除有序集合成员: zset_key removed member1 = %v\n", removed > 0)
}

// 演示自定义命令
func demoCustomCommands(ctx context.Context, client cache.Client) {
	fmt.Println("\n--- 自定义命令 ---")

	// 执行自定义命令
	result, err := client.Do(ctx, "ping")
	if err != nil {
		log.Printf("自定义命令失败: %v", err)
		return
	}
	fmt.Printf("执行PING命令: result = %v\n", result)

	// 执行INFO命令
	result, err = client.Do(ctx, "info", "server")
	if err != nil {
		log.Printf("INFO命令失败: %v", err)
		return
	}
	fmt.Printf("执行INFO命令: \n%v\n", result)
}
