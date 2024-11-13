Env NOTEDOWN_DIR "features/addtask/final"
Env TEST_DATE "2020-01-01"

Source "features/common.tape"
Output "features/addtask/demo.gif"

Type "a"
Space
Tab
Type "This is a basic task"
Sleep 1
Enter
Sleep 1

Type "a"
Type "b"
Tab
Type "This is a blocked task with a high priority p:1"
Sleep 1
Enter
Sleep 1  

Type "a"
Type "x"
Tab
Type "This is a completed task"
Sleep 1
Enter
Sleep 1  

Type "a"
Type "/"
Tab
Type "This is an in progress task that repeats every:day"
Sleep 1
Enter

Sleep 5
