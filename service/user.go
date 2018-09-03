package service

import (
	"apiserver/model"
	"sync"
	"apiserver/util"
	"fmt"
)

/*做具体的查询处理*/
func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
	infos := make([]*model.UserInfo, 0) //存储返回结果切片

	//数据库查询
	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}
	//构建id集合的切片，保存查询的顺序
	ids := []uint64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	wg := sync.WaitGroup{}
	//每个goroutine用自己的锁，以及map映射
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.UserInfo, len(users)),//记录顺序
	}

	//channel，用于循环并发，routine错误
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	//循环users，写入UserInfo封装
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserModel) {
			defer wg.Done()
			//生成Id
			shortId, err := util.GenShortId()
			if err != nil {
				errChan <- err
				return
			}

			//对UserInfo加锁，同一时间只能有一个在操作，封装
			//更新同一个变量为了保证数据一致性
			userList.Lock.Lock()
			defer userList.Lock.Unlock()
			userList.IdMap[u.Id] = &model.UserInfo{ //使用了goroutine,顺序就会乱
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan: //处理错误
		return nil, count, err
	}

	//返回infos
	for _, id := range ids { //便于重新记录顺序
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
