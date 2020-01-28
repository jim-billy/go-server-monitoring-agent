package main
import (
	"fmt"
	"com/coder/executor"
)

func testExecutor(){
	//var exec executor.Executor
	exec := new(executor.Executor)
	
	exec.SetCommand("/home/test/git/me_agent/product_package/scripts/script.sh")
	exec.SetCommandArgs([]string{"cpu_util"})
	exec.SetTimeout(10)
	exec.Execute()
	fmt.Println("===================== Is success : ",exec.IsSuccess(),", Execution time : ",exec.GetExecutionTime(),", Output ",exec.GetOutput(),", Error : ",exec.GetError())
	
	exec1 := new(executor.Executor)
	args := []string{"-c 4", "-i 1", "8.8.8.8"}
	exec1.SetCommand("ping")
	exec1.SetCommandArgs(args)
	exec1.SetTimeout(12)
	exec1.Execute()
	fmt.Println("===================== Is success : ",exec1.IsSuccess(),", Execution time : ",exec1.GetExecutionTime(),", Output ",exec1.GetOutput(),", Error : ",exec1.GetError())
	
}

func main(){
	testExecutor()
}