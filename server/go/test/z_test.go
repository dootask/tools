package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	dootask "github.com/dootask/tools/server/go"
)

// 测试配置常量
const (
	testToken  = "YIG8ANC8q2ROQF91r8Pe6-53rIG3oCxcqQN-mMdZpQKe7mKwNqIHenDNqbDDdyQIdo9w2KdveEpF1NaH-5Nfmv0dBr9TkjJ7KFMkfEUL11wOjyId0nuoSJaAliRz8d5z"
	testServer = "http://127.0.0.1:2222"

	// 测试用的用户ID和对话ID
	testUserID    = 3
	testUserID2   = 33
	testDialogID  = 2367
	testDialogID2 = 2889

	// 时间格式常量
	timeFormat = "2006-01-02 15:04:05"
)

// formatJSON 格式化 JSON 输出，用于测试日志
func formatJSON(v interface{}) string {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("JSON序列化失败: %v", err)
	}
	return string(jsonBytes)
}

// setupTestClient 创建测试客户端
func setupTestClient() *dootask.Client {
	return dootask.NewClient(testToken, dootask.WithServer(testServer), dootask.WithTimeout(30*time.Second))
}

// ============================================================================
// 用户相关测试
// ============================================================================

func TestUserAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("获取用户信息", func(t *testing.T) {
		user, err := client.GetUserInfo()
		if err != nil {
			t.Fatalf("获取用户信息失败: %v", err)
		}

		if user == nil {
			t.Fatal("用户信息为空")
		}

		t.Logf("用户信息:\n%s", formatJSON(user))
	})

	t.Run("获取用户基础信息", func(t *testing.T) {
		userBasic, err := client.GetUserBasic(testUserID)
		if err != nil {
			t.Fatalf("获取用户基础信息失败: %v", err)
		}

		if userBasic == nil {
			t.Fatal("用户基础信息为空")
		}

		t.Logf("用户基础信息:\n%s", formatJSON(userBasic))
	})

	t.Run("获取多个用户基础信息", func(t *testing.T) {
		userIDs := []int{testUserID, testUserID2}
		userBasics, err := client.GetUsersBasic(userIDs)
		if err != nil {
			t.Fatalf("获取多个用户基础信息失败: %v", err)
		}

		if len(userBasics) == 0 {
			t.Fatal("用户基础信息列表为空")
		}

		t.Logf("多个用户基础信息:\n%s", formatJSON(userBasics))
	})
}

// ============================================================================
// 消息相关测试
// ============================================================================

func TestMessageAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("发送消息", func(t *testing.T) {
		req := dootask.SendMessageRequest{
			DialogID: testDialogID,
			Text:     "Hello, world! 发送指定会话消息 " + time.Now().Format(timeFormat),
		}

		err := client.SendMessage(req)
		if err != nil {
			t.Fatalf("发送消息失败: %v", err)
		}

		t.Log("消息发送成功")
	})

	t.Run("发送消息到用户", func(t *testing.T) {
		req := dootask.SendMessageToUserRequest{
			UserID: testUserID,
			Text:   "Hello, world! 发送消息到用户 " + time.Now().Format(timeFormat),
		}

		err := client.SendMessageToUser(req)
		if err != nil {
			t.Fatalf("发送消息到用户失败: %v", err)
		}

		t.Log("消息发送到用户成功")
	})

	t.Run("发送机器人消息", func(t *testing.T) {
		req := dootask.SendBotMessageRequest{
			UserID: testUserID,
			Text:   "Hello, world! 发送机器人消息 " + time.Now().Format(timeFormat),
		}

		err := client.SendBotMessage(req)
		if err != nil {
			t.Fatalf("发送机器人消息失败: %v", err)
		}

		t.Log("机器人消息发送成功")
	})

	t.Run("发送匿名消息", func(t *testing.T) {
		req := dootask.SendAnonymousMessageRequest{
			UserID: testUserID,
			Text:   "Hello, world! 发送匿名消息 " + time.Now().Format(timeFormat),
		}

		err := client.SendAnonymousMessage(req)
		if err != nil {
			t.Fatalf("发送匿名消息失败: %v", err)
		}

		t.Log("匿名消息发送成功")
	})
}

// ============================================================================
// 对话相关测试
// ============================================================================

func TestDialogAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("获取对话列表", func(t *testing.T) {
		req := dootask.TimeRangeRequest{}
		dialogs, err := client.GetDialogList(req)
		if err != nil {
			t.Fatalf("获取对话列表失败: %v", err)
		}

		if dialogs == nil {
			t.Fatal("对话列表为空")
		}

		t.Logf("对话列表:\n%s", formatJSON(dialogs))
	})

	t.Run("获取单个对话信息", func(t *testing.T) {
		req := dootask.GetDialogRequest{DialogID: testDialogID2}
		dialog, err := client.GetDialogOne(req)
		if err != nil {
			t.Fatalf("获取单个对话信息失败: %v", err)
		}

		if dialog == nil {
			t.Fatal("对话信息为空")
		}

		t.Logf("对话信息:\n%s", formatJSON(dialog))
	})

	t.Run("获取对话成员", func(t *testing.T) {
		req := dootask.GetDialogUserRequest{DialogID: testDialogID2}
		dialogMembers, err := client.GetDialogUser(req)
		if err != nil {
			t.Fatalf("获取对话成员失败: %v", err)
		}

		if dialogMembers == nil {
			t.Fatal("对话成员信息为空")
		}

		t.Logf("对话成员:\n%s", formatJSON(dialogMembers))
	})
}

// ============================================================================
// 群组管理测试
// ============================================================================

func TestGroupAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("创建群组", func(t *testing.T) {
		req := dootask.CreateGroupRequest{
			ChatName: "测试群组" + time.Now().Format(timeFormat),
			UserIDs:  []int{testUserID, testUserID2},
		}

		group, err := client.CreateGroup(req)
		if err != nil {
			t.Fatalf("创建群组失败: %v", err)
		}

		t.Logf("群组创建成功: %s", formatJSON(group))
	})

	t.Run("修改群组", func(t *testing.T) {
		req := dootask.EditGroupRequest{
			DialogID: 27211,
			ChatName: "测试修改群组" + time.Now().Format(timeFormat),
		}

		err := client.EditGroup(req)
		if err != nil {
			t.Fatalf("修改群组失败: %v", err)
		}

		t.Log("群组修改成功")
	})

	t.Run("添加群组成员", func(t *testing.T) {
		req := dootask.AddGroupUserRequest{
			DialogID: 27211,
			UserIDs:  []int{34},
		}

		err := client.AddGroupUser(req)
		if err != nil {
			t.Fatalf("添加群组成员失败: %v", err)
		}

		t.Log("群组成员添加成功")
	})

	t.Run("移除群组成员", func(t *testing.T) {
		req := dootask.RemoveGroupUserRequest{
			DialogID: 27211,
			UserIDs:  []int{34},
		}

		err := client.RemoveGroupUser(req)
		if err != nil {
			t.Fatalf("移除群组成员失败: %v", err)
		}

		t.Log("群组成员移除成功")
	})

	t.Run("转让群组", func(t *testing.T) {
		req := dootask.TransferGroupRequest{
			DialogID: 27211,
			UserID:   33,
		}

		err := client.TransferGroup(req)
		if err != nil {
			t.Fatalf("转让群组失败: %v", err)
		}

		t.Log("群组转让成功")
	})
}

// ============================================================================
// 项目管理测试
// ============================================================================

func TestProjectAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("项目管理完整流程", func(t *testing.T) {
		// 1. 获取项目列表
		t.Log("--- 测试获取项目列表 ---")
		params := dootask.GetProjectListRequest{
			Type:          "all",
			Archived:      "no",
			GetStatistics: "yes",
			Page:          1,
			PageSize:      10,
		}

		response, err := client.GetProjectList(params)
		if err != nil {
			t.Fatalf("获取项目列表失败: %v", err)
		}

		t.Logf("✓ 项目列表获取成功，总数: %d", response.Total)
		for i, project := range response.Data {
			if i < 3 { // 只显示前3个项目
				t.Logf("  - 项目: %s (ID: %d)", project.Name, project.ID)
			}
		}

		// 2. 创建项目
		t.Log("--- 测试创建项目 ---")
		createParams := dootask.CreateProjectRequest{
			Name:    "测试项目-" + time.Now().Format("20060102150405"),
			Desc:    "这是一个测试项目",
			Columns: "待办,进行中,已完成",
			Flow:    "close",
		}

		project, err := client.CreateProject(createParams)
		if err != nil {
			t.Fatalf("创建项目失败: %v", err)
		}

		t.Logf("✓ 项目创建成功: %s (ID: %d)", project.Name, project.ID)
		projectID := project.ID

		// 3. 获取项目信息
		t.Log("--- 测试获取项目信息 ---")
		getParams := dootask.GetProjectRequest{
			ProjectID: projectID,
		}

		project, err = client.GetProject(getParams)
		if err != nil {
			t.Fatalf("获取项目信息失败: %v", err)
		}

		t.Logf("✓ 项目信息获取成功: %s", project.Name)

		// 4. 更新项目
		t.Log("--- 测试更新项目 ---")
		updateParams := dootask.UpdateProjectRequest{
			ProjectID: projectID,
			Name:      "更新后的测试项目",
			Desc:      "这是一个更新后的测试项目描述",
		}

		project, err = client.UpdateProject(updateParams)
		if err != nil {
			t.Fatalf("更新项目失败: %v", err)
		}

		t.Logf("✓ 项目更新成功: %s", project.Name)

		// 5. 创建任务列表
		t.Log("--- 测试创建任务列表 ---")
		columnParams := dootask.CreateColumnRequest{
			ProjectID: projectID,
			Name:      "测试列表",
		}

		column, err := client.CreateColumn(columnParams)
		if err != nil {
			t.Fatalf("创建任务列表失败: %v", err)
		}

		t.Logf("✓ 任务列表创建成功: %s (ID: %d)", column.Name, column.ID)
		columnID := column.ID

		// 6. 获取任务列表
		t.Log("--- 测试获取任务列表 ---")
		getColumnParams := dootask.GetColumnListRequest{
			ProjectID: projectID,
			Page:      1,
			PageSize:  10,
		}

		columnResponse, err := client.GetColumnList(getColumnParams)
		if err != nil {
			t.Fatalf("获取任务列表失败: %v", err)
		}

		t.Logf("✓ 任务列表获取成功，总数: %d", columnResponse.Total)

		// 7. 更新任务列表
		t.Log("--- 测试更新任务列表 ---")
		updateColumnParams := dootask.UpdateColumnRequest{
			ColumnID: columnID,
			Name:     "更新后的测试列表",
			Color:    "#FF0000",
		}

		column, err = client.UpdateColumn(updateColumnParams)
		if err != nil {
			t.Fatalf("更新任务列表失败: %v", err)
		}

		t.Logf("✓ 任务列表更新成功: %s", column.Name)

		// 8. 创建任务
		t.Log("--- 测试创建任务 ---")
		taskParams := dootask.CreateTaskRequest{
			ProjectID: projectID,
			ColumnID:  columnID,
			Name:      "测试任务-" + time.Now().Format("20060102150405"),
			Content:   "这是一个测试任务的内容",
			Times:     []string{"2024-01-01 09:00:00", "2024-01-01 18:00:00"},
		}

		task, err := client.CreateTask(taskParams)
		if err != nil {
			t.Fatalf("创建任务失败: %v", err)
		}

		t.Logf("✓ 任务创建成功: %s (ID: %d)", task.Name, task.ID)
		taskID := task.ID

		// 9. 获取任务信息
		t.Log("--- 测试获取任务信息 ---")
		getTaskParams := dootask.GetTaskRequest{
			TaskID:   taskID,
			Archived: "no",
		}

		task, err = client.GetTask(getTaskParams)
		if err != nil {
			t.Fatalf("获取任务信息失败: %v", err)
		}

		t.Logf("✓ 任务信息获取成功: %s", task.Name)

		// 10. 获取任务列表
		t.Log("--- 测试获取任务列表 ---")
		getTaskListParams := dootask.GetTaskListRequest{
			ProjectID: projectID,
			Archived:  "no",
			Page:      1,
			PageSize:  10,
		}

		taskResponse, err := client.GetTaskList(getTaskListParams)
		if err != nil {
			t.Fatalf("获取任务列表失败: %v", err)
		}

		t.Logf("✓ 任务列表获取成功，总数: %d", taskResponse.Total)

		// 11. 获取任务内容
		t.Log("--- 测试获取任务内容 ---")
		getContentParams := dootask.GetTaskContentRequest{
			TaskID: taskID,
		}

		_, err = client.GetTaskContent(getContentParams)
		if err != nil {
			t.Fatalf("获取任务内容失败: %v", err)
		}

		t.Log("✓ 任务内容获取成功")

		// 12. 获取任务文件
		t.Log("--- 测试获取任务文件 ---")
		getFilesParams := dootask.GetTaskFilesRequest{
			TaskID: taskID,
		}

		files, err := client.GetTaskFiles(getFilesParams)
		if err != nil {
			t.Fatalf("获取任务文件失败: %v", err)
		}

		t.Logf("✓ 任务文件获取成功，文件数: %d", len(files))

		// 13. 创建子任务
		t.Log("--- 测试创建子任务 ---")
		subTaskParams := dootask.CreateSubTaskRequest{
			TaskID: taskID,
			Name:   "测试子任务",
		}

		subTask, err := client.CreateSubTask(subTaskParams)
		if err != nil {
			t.Fatalf("创建子任务失败: %v", err)
		}

		t.Logf("✓ 子任务创建成功: %s (ID: %d)", subTask.Name, subTask.ID)

		// 14. 更新任务
		t.Log("--- 测试更新任务 ---")
		updateTaskParams := dootask.UpdateTaskRequest{
			TaskID:  taskID,
			Name:    "更新后的测试任务",
			Content: "这是更新后的任务内容",
			Color:   "#00FF00",
		}

		task, err = client.UpdateTask(updateTaskParams)
		if err != nil {
			t.Logf("更新任务失败: %v", err)
		} else {
			t.Logf("✓ 任务更新成功: %s", task.Name)
		}

		// 15. 创建任务对话
		t.Log("--- 测试创建任务对话 ---")
		dialogParams := dootask.CreateTaskDialogRequest{
			TaskID: taskID,
		}

		dialog, err := client.CreateTaskDialog(dialogParams)
		if err != nil {
			t.Fatalf("创建任务对话失败: %v", err)
		}

		t.Logf("✓ 任务对话创建成功: 任务ID: %d, 对话ID: %d", dialog.ID, dialog.DialogID)

		// 清理：删除任务列表
		t.Log("--- 清理：删除任务列表 ---")
		err = client.DeleteColumn(columnID)
		if err != nil {
			t.Fatalf("删除任务列表失败: %v", err)
		}

		t.Logf("✓ 任务列表删除成功 (ID: %d)", columnID)

		// 清理：删除项目
		t.Log("--- 清理：删除项目 ---")
		err = client.DeleteProject(projectID)
		if err != nil {
			t.Fatalf("删除项目失败: %v", err)
		}

		t.Logf("✓ 项目删除成功 (ID: %d)", projectID)

		t.Log("=== 项目管理测试完成 ===")
	})
}

// ============================================================================
// 基础功能测试
// ============================================================================

func TestBasicProjectAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("获取项目列表", func(t *testing.T) {
		req := dootask.GetProjectListRequest{}
		projects, err := client.GetProjectList(req)
		if err != nil {
			t.Fatalf("获取项目列表失败: %v", err)
		}

		t.Logf("项目列表: %s", formatJSON(projects))
	})
}

// ============================================================================
// 消息相关测试
// ============================================================================

func TestMessageListAPI(t *testing.T) {
	client := setupTestClient()

	t.Run("获取消息列表", func(t *testing.T) {
		req := dootask.GetMessageListRequest{
			DialogID: 5176,
			Take:     10,
		}
		messages, err := client.GetMessageList(req)
		if err != nil {
			t.Fatalf("获取消息列表失败: %v", err)
		}

		t.Logf("消息列表: %s", formatJSON(messages))
	})
}

// ============================================================================
// 系统设置测试
// ============================================================================

func TestGetSystemSettings(t *testing.T) {
	client := setupTestClient()

	t.Run("获取系统设置", func(t *testing.T) {
		settings, err := client.GetSystemSettings()
		if err != nil {
			t.Fatalf("获取系统设置失败: %v", err)
		}
		t.Logf("系统设置: %s", formatJSON(settings))
	})
}

func TestGetVersion(t *testing.T) {
	client := setupTestClient()

	t.Run("获取版本信息", func(t *testing.T) {
		version, err := client.GetVersion()
		if err != nil {
			t.Fatalf("获取版本信息失败: %v", err)
		}
		t.Logf("版本信息: %s", formatJSON(version))
	})
}
