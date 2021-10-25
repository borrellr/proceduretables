package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "container/list"
)

const (
   GROUP_STACK_SIZE = 50
   NO_ERROR = 0
   EXIT_NOW = 2001
)

type Func func(pc *prcontrol) int

type prcontrol struct {
    current_question *list.Element
    current_group *list.List
    response *string
    errstat int
    errormess map[int]errormesstype
}

type Question struct {
    text *string
    response *string
    validate, doit, set Func
}

type group_stack_type struct {
    current_group *list.List
    current_question *list.Element
}

type errormesstype struct {
//    errstat int
    message *string
    build func(msg *string) int
}

var group_stack = list.New()

func prompter(pc *prcontrol) int {
    pc.current_question = pc.current_group.Front()
//    pc.group_stack_ptr = 0

    for {
        cq := pc.current_question.Value.(Question)
	pc.errstat = NO_ERROR

	display_current_question(pc)

	get_response(pc)

	if pc.response == nil {
            continue
        }

	if pc.errstat = cq.validate(pc); pc.errstat == NO_ERROR {
	    if pc.errstat = cq.doit(pc); pc.errstat > NO_ERROR {
	       if (pc.errstat == EXIT_NOW) {
		   return NO_ERROR
               }
            }
	}


	if cq.response != nil {
	    cq.response = pc.response
        }

	cq.set(pc)

	if pc.errstat > NO_ERROR {
	    handle_error(pc.errstat, pc.errormess)
        }
    }
}

func display_current_question(pc *prcontrol) {
    cq := pc.current_question.Value.(Question)
    if cq.text != nil {
       fmt.Printf("\r\n%s\r\n", *cq.text)
       fmt.Printf("---> ")
    }
}

func get_response(pc *prcontrol) {
    cq := pc.current_question.Value.(Question)
    if cq.text != nil {
       reader := bufio.NewReader(os.Stdin)
       str, _ := reader.ReadString('\n')
       str = strings.TrimSuffix(str, "\r\n")
       str = strings.ToLower(str)
       pc.response = &str
    }
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
     pc.current_question = pc.current_question.Next()
     return NO_ERROR
}

func pop_group(pc *prcontrol) int {
     e := group_stack.Front()
     if e != nil {
         gsp := e.Value.(group_stack_type)
         pc.current_group = gsp.current_group
         pc.current_question = gsp.current_question
	 group_stack.Remove(e)
     }
     return NO_ERROR
}

func push_current_group(pc *prcontrol) int {
     gsp := group_stack_type {
	     current_group : pc.current_group,
	     current_question : pc.current_question,
     }
     group_stack.PushFront(gsp)
     return NO_ERROR
}

func start_group(newgroup *list.List, pc *prcontrol) int {
     push_current_group(pc)
     pc.current_group = newgroup
     pc.current_question = newgroup.Front()
     return NO_ERROR
}

func restart_group(pc *prcontrol) int {
     pc.current_question = pc.current_group.Front()
     return NO_ERROR
}

func end_group(pc *prcontrol) int {
     pop_group(pc)
     pc.current_question = pc.current_question.Next()
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

func prtgrpstk() {
     fmt.Println("Checking the group stack...")
     for e := group_stack.Front() ; e != nil ; e = e.Next() {
	 item := e.Value.(group_stack_type)
	 fmt.Println(item)
     }
}
