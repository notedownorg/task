Env NOTEDOWN_DIR "features/addtask/final"
Source "features/common.tape"
Output "features/addtask/scenario.gif"

Type "a"
Space
Tab
Type "This is a basic task"
Sleep 0.5  
Enter
Sleep 0.5  

Type "a"
Type "b"
Tab
Type "This is a blocked task with a high priority p:1"
Sleep 0.5
Enter
Sleep 0.5  

Type "a"
Type "x"
Tab
Type "This is a completed task"
Sleep 0.5
Enter
Sleep 0.5  

Type "a"
Type "/"
Tab
Type "This is an in progress task that repeats every:day"
Sleep 0.5
Enter

Sleep 5
