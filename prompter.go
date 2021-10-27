package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "container/list"
)

const (
   NO_ERROR = 0
   EXIT_NOW = 2001
)

type Func func(pc *prcontrol) int

type prcontrol struct {
    current_question *list.Element
    current_group *list.List
    group_stack *list.List
    response *string
    errstat int
    errormess map[int]errormesstype
}

func NewPrControl() *prcontrol {
     var pc prcontrol
     pc.group_stack = list.New()
     return &pc
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
    message *string
    build func(msg *string) int
}

//*****************************
//  prcontrol methods
//*****************************

func (pc *prcontrol) prompter() int {
    pc.current_question = pc.current_group.Front()

    for {
        cq := pc.current_question.Value.(Question)
	pc.errstat = NO_ERROR

	pc.display_current_question()

	pc.get_response()

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

func (pc *prcontrol) display_current_question() {
    cq := pc.current_question.Value.(Question)
    if cq.text != nil {
       fmt.Printf("\r\n%s\r\n", *cq.text)
       fmt.Printf("---> ")
    }
}

func (pc *prcontrol) get_response() {
    cq := pc.current_question.Value.(Question)
    if cq.text != nil {
       reader := bufio.NewReader(os.Stdin)
       str, _ := reader.ReadString('\n')
       str = strings.TrimSuffix(str, "\r\n")
       str = strings.ToLower(str)
       pc.response = &str
    }
}

func (pc *prcontrol) next_question() {
     pc.current_question = pc.current_question.Next()
}

func (pc *prcontrol) pop_group() {
     e := pc.group_stack.Front()
     if e != nil {
         gsp := e.Value.(group_stack_type)
         pc.current_group = gsp.current_group
         pc.current_question = gsp.current_question
	 pc.group_stack.Remove(e)
     }
}

func (pc *prcontrol) push_current_group() {
     gsp := group_stack_type {
	     current_group : pc.current_group,
	     current_question : pc.current_question,
     }
     pc.group_stack.PushFront(gsp)
}

func (pc *prcontrol) end_group() {
     pc.pop_group()
     pc.current_question = pc.current_question.Next()
}

func (pc *prcontrol) restart_group() {
     pc.current_question = pc.current_group.Front()
}

func (pc *prcontrol) start_group(newgroup *list.List) {
     pc.push_current_group()
     pc.current_group = newgroup
     pc.current_question = newgroup.Front()
}

//*****************************
//  error handle function
//*****************************

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

//****************************
//  Flow control routines
//****************************

func prgexit(pc *prcontrol) int {
     return EXIT_NOW
}

func no_op(pc *prcontrol) int {
     return NO_ERROR
}

func checkerror_end_group(pc *prcontrol) int {
     if pc.errstat > NO_ERROR {
	return NO_ERROR
     }
     pc.end_group()
     return NO_ERROR
}

func checkerror_next_question(pc *prcontrol) int {
     if pc.errstat > 0 {
        return NO_ERROR
     }
     pc.next_question()
     return NO_ERROR
}
