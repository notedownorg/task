Require go

Set Width 1200
Set Height 600
Set FontSize 12
Set FontFamily "JetBrainsMono Nerd Font"

Env NOTEDOWN_DIR "testdata/tmp/addtask"
Env TEST_DATE "2020-01-01"

Hide
Type "go run ../main.go"
Enter
Sleep 3
Show

Type "a"
Space
Tab
Type "This is a task"
Enter

Type "a"
Type "b"
Tab
Type "This is a blocked task p:1"
Enter

Type "a"
Type "x"
Tab
Type "This is a completed task "
Enter

Type "a"
Type "/"
Tab
Type "This is an in progress task every:day"
Enter

Sleep 5
