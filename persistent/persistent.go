/*************************  mongodb backend  (mongodb接口测试程序)     ***********************/

// 搭建docker mongo环境：
//      docker run -d --network host mongo:latest
//      docker run --name  my-mongo  --net=host -p 27017:27017  -d mongo
//      docker exec -it 52be7c6761c9 bash

// 下载官方go mongo-driver:
//     1. go get github.com/mongodb/mongo-go-driver
//     2. 将github.com/mongodb/mongo-go-driver/重命名为go.mongodb.org/mongo-driver/

package persistent

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	Name string
	Age  int
}

func Test() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 指定获取要操作的数据集
	collection := client.Database("test").Collection("student")

	// 插入文档
	if true {
		s1 := Student{"小红", 12}
		s2 := Student{"小兰", 10}
		s3 := Student{"小黄", 11}
		// 插入一条文档记录
		insertResult, err := collection.InsertOne(context.TODO(), s1) //interface{}是空接口，它可以承接任意类型，然后进行推导
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		// 插入多条文档记录
		students := []interface{}{s2, s3} //[]interface{}是空接口切片，它可以承接任意类型切片，然后进行推导
		insertManyResult, err := collection.InsertMany(context.TODO(), students)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	}

	// 更新文档
	if true {
		filter := bson.D{{"name", "小兰"}}
		update := bson.D{
			{"$inc", bson.D{ //$inc操作符将一个字段的值增加或者减少
				{"age", 1},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	// 查找文档
	if true {
		filter := bson.D{{"name", "小兰"}}
		var result Student
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found a single document: %+v\n", result)
	}

	// 查询多个文档
	if true {
		findOptions := options.Find()
		findOptions.SetLimit(2)

		var results []*Student //指针切片

		// 返回用来遍历的游标cursor
		cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {
			var elem Student
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		cur.Close(context.TODO())
		// fmt.Printf("Found multiple documents (array of pointers): %#v\n", results)
		for k, v := range results {
			fmt.Printf("key:%v value:%v value type:%T\n", k, v, v)
		}
	}

	// 删除单个文档
	if true {
		deleteResult1, err := collection.DeleteOne(context.TODO(), bson.D{{"name", "小黄"}})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)
	}

	// 删除所有文档
	if true {
		deleteResult2, err := collection.DeleteMany(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)
	}

	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
