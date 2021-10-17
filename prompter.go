package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
)

const (
   GROUP_STACK_SIZE = 50
   NO_ERROR = 0
   EXIT_NOW = 2001
)

type Func func(pc *prcontrol) int

type prcontrol struct {
    current_question int
    current_group []question_type
    group_stack_ptr int
    response *string
    errstat int
    errormess map[int]errormesstype
}

type question_type struct {
    text *string
    response *string
    validate, doit, set Func
}

type group_stack_type struct {
    current_group []question_type
    current_question int
}

type errormesstype struct {
//    errstat int
    message *string
    build func(msg *string) int
}

var group_stack = make([]group_stack_type, GROUP_STACK_SIZE)

func prompter(pc *prcontrol) int {
//    var errstat int
    pc.current_question = 0
    pc.group_stack_ptr = 0

    for {
	pc.errstat = NO_ERROR

	display_current_question(pc)

	get_response(pc)

	if pc.response == nil {
            continue
        }

	if pc.errstat = pc.current_group[pc.current_question].validate(pc); pc.errstat == NO_ERROR {
	    if pc.errstat = pc.current_group[pc.current_question].doit(pc); pc.errstat > NO_ERROR {
	       if (pc.errstat == EXIT_NOW) {
		   return NO_ERROR
               }
            }
	}

	if pc.current_group[pc.current_question].response != nil {
	    pc.current_group[pc.current_question].response = pc.response
        }

	pc.current_group[pc.current_question].set(pc)

	if pc.current_group[pc.current_question].text == nil {
	    return NO_ERROR
	}

	if pc.errstat > NO_ERROR {
	    handle_error(pc.errstat, pc.errormess)
        }
    }
}

func display_current_question(pc *prcontrol) {
    fmt.Printf("\r\n%s\r\n", *pc.current_group[pc.current_question].text)
    fmt.Printf("---> ")
}

func get_response(pc *prcontrol) {
    reader := bufio.NewReader(os.Stdin)
    str, _ := reader.ReadString('\n')
    str = strings.TrimSuffix(str, "\n")
    str = strings.ToLower(str)
    pc.response = &str
}

func handle_error(errstat int, errormess map[int]errormesstype) int {
    var message errormesstype

    if (errstat > -1) {
       message = errormess[errstat]
       if message.build != nil {
	   message.build(message.message)
       } else {
          fmt.Printf("%s Error %d.\n\r", *message.message, errstat)
       }
    }
//    fmt.Sprintf("\n\r%s\n\r", message)
    return NO_ERROR
}

/*****************************
    Flow control routines
*****************************/
func prgexit(pc *prcontrol) int {
     return EXIT_NOW
}

func no_op(pc *prcontrol) int {
     return NO_ERROR
}

func next_question(pc *prcontrol) int {
     pc.current_question++
     return NO_ERROR
}

func pop_group(pc *prcontrol) int {
     pc.group_stack_ptr--
     pc.current_group = group_stack[pc.group_stack_ptr].current_group
     pc.current_question = group_stack[pc.group_stack_ptr].current_question
     return NO_ERROR
}

func push_current_group(pc *prcontrol) int {
     group_stack[pc.group_stack_ptr].current_group = pc.current_group
     group_stack[pc.group_stack_ptr].current_question = pc.current_question
     pc.group_stack_ptr++
     return NO_ERROR
}

func start_group(newgroup []question_type, pc *prcontrol) int {
     push_current_group(pc)
     pc.current_group = newgroup
     pc.current_question = 0
     return NO_ERROR
}

func restart_group(pc *prcontrol) int {
     pc.current_question = 0
     return NO_ERROR
}

func end_group(pc *prcontrol) int {
     pop_group(pc)
     pc.current_question++
     return NO_ERROR
}

func checkerror_end_group(pc *prcontrol) int {
     if pc.errstat > NO_ERROR {
	return NO_ERROR
     }
     end_group(pc)
     return NO_ERROR
}

func checkerror_next_question(pc *prcontrol) int {
     if pc.errstat > 0 {
        return NO_ERROR
     }
     next_question(pc)
     return NO_ERROR
}
