package main

import (
    "strings"
    "strconv"
    "fmt"
    "os"
    "container/list"
)

const (
   ENTER_S_OR_R = iota + 1
   ENTER_Y_OR_N
   START_ACCOUNT_LARGER
   BAD_PARAM_NAME
   BAD_ACCOUNT_NUMBER
   ENTER_P_S_OR_D
   FILE_EXISTS
   BAD_PARM_NAME
)

//**********************************
//  The report parameter variables
//**********************************

var (
	report_destination string
	dest_filename string
	single_or_range string
	start_account string
	end_account string
	account_number int
	display_parmname string
	include_overshort string
	endq Question
)

//***********************
// Report to printer, screen or disk routines
//***********************

var report_filename = list.New()

func init() {
    msg := "What is the name of the disk file?"
    q := Question {
	    text : &msg,
	    response: &dest_filename,
	    validate: filename_val,
	    doit : no_op,
	    set : checkerror_end_group }
    report_filename.PushBack(q)
}

func filename_val(pc *prcontrol) int {
     ret_code := NO_ERROR
     fname := *pc.response

     file, err := os.Open(fname)
     if err == nil {
	file.Close()
	ret_code = FILE_EXISTS
     }
     return ret_code
}

func reportdest_val(pc *prcontrol) int {
     resp := *pc.response
     var ret_code int

     if strings.HasPrefix(resp, "p") || strings.HasPrefix(resp, "s") || strings.HasPrefix(resp, "d") {
	ret_code = NO_ERROR
     } else {
	ret_code = ENTER_P_S_OR_D
     }
     return ret_code
}

func reportdest_set(pc *prcontrol) int {
     resp := *pc.response

     if strings.HasPrefix(resp, "d") {
	start_group(report_filename, pc)
     } else if strings.HasPrefix(resp, "p") || strings.HasPrefix(resp, "s") {
	next_question(pc)
     }
     return NO_ERROR
}

//**************************
// Account routines
//**************************

var account_range = list.New()

func init() {
     var str [2]string

     str[0] = "Enter the starting account."
     str[1] = "Enter the ending account."

     q := Question {
	 text : &str[0],
	 response : &start_account,
	 validate : account_val,
	 doit : no_op,
	 set : checkerror_next_question }

     q1 := Question {
	 text : &str[1],
	 response : &end_account,
	 validate : end_account_val,
	 doit : no_op,
	 set : end_account_set }

     account_range.PushBack(q)
     account_range.PushBack(q1)
}

var account = list.New()

func init() {
     str := "Enter the Account."

     q := Question {
	 text : &str,
	 response : &start_account,
	 validate : account_val,
	 doit : save_account_doit,
	 set : checkerror_end_group }

     account.PushBack(q)
}

//*****************************
//   Account routines
//*****************************

func account_or_range_val(pc *prcontrol) int {
     resp := *pc.response
     if (strings.HasPrefix(resp, "r") || strings.HasPrefix(resp, "s")) {
	 return NO_ERROR
     } else {
	 return ENTER_S_OR_R
     }
}

func account_or_range_set(pc *prcontrol) int {
     if pc.errstat == NO_ERROR {
        if strings.HasPrefix(*pc.response, "r") {
	    start_group(account_range, pc)
        } else if strings.HasPrefix(*pc.response, "s") {
	    start_group(account, pc)
        }
     }
     return NO_ERROR
}

func save_account_doit(pc *prcontrol) int {
     account_number, _ = strconv.Atoi(*pc.response)
     return NO_ERROR
}

func account_val(pc *prcontrol) int {
     if resp, _ := strconv.Atoi(*pc.response); resp > 99 && resp < 1001 {
	return NO_ERROR
     } else {
	return BAD_ACCOUNT_NUMBER
     }
}

func end_account_val(pc *prcontrol) int {
    st, _ := strconv.Atoi(start_account)

    if errstat := account_val(pc); errstat > NO_ERROR {
       return errstat
    }

    if  resp, _ := strconv.Atoi(*pc.response); st > resp {
       return START_ACCOUNT_LARGER
    }
    return NO_ERROR
}

func end_account_set(pc *prcontrol) int {
    switch(pc.errstat) {
	case NO_ERROR:
	     end_group(pc)
	case START_ACCOUNT_LARGER:
	     restart_group(pc)
        case BAD_ACCOUNT_NUMBER:
    }
    return NO_ERROR
}

//********************************
// Get display parameters routines
//********************************

// In a "real" system, this table
// would probably be stored in a file
// and parmname_val would check to see
// if the name entered is in this file.

var legal_parmnames = [4]string{"default", "daily", "weekly", "yearly" }

func parmname_val(pc *prcontrol) int {
     resp := *pc.response

     for _, lp := range legal_parmnames {
	 if lp == resp {
	    return NO_ERROR
	 }
     }
     return BAD_PARM_NAME
}

func bld_bad_parmname(message *string) int {
     fmt.Fprintf(os.Stdout, "%s %s.\n", *message, legal_parmnames)
     return NO_ERROR
}

//***************************
//   yesno validation
//***************************

func yesno_val(pc *prcontrol) int {
	resp := *pc.response
     if (strings.HasPrefix(resp, "y") || strings.HasPrefix(resp, "n")) {
	 return NO_ERROR
     } else {
	 return ENTER_Y_OR_N
     }
}


//**************************************
// Main question array procedure table
//**************************************

var account_parms = list.New()

func init () {
     var account_parms_text = make([]string, 4)
     account_parms_text[0] = "Do you want this report for a single account or a range of accounts? (S or R)"
     account_parms_text[1] = "Enter the name of the display parameter record."
     account_parms_text[2] = "Do you want to include the Over/Short Report? (Y/N)"
     account_parms_text[3] = "Do you want this report on the printer, screen, or saved to disk?(P,S or D)"

     q := Question {
	 text : &account_parms_text[0],
	 response : &single_or_range,
	 validate : account_or_range_val,
         doit : no_op,
         set : account_or_range_set }

     q1 := Question {
	 text : &account_parms_text[1],
	 response : &display_parmname,
	 validate : parmname_val,
         doit : no_op,
         set : checkerror_next_question }

     q2 := Question {
	 text : &account_parms_text[2],
	 response : nil,
	 validate : yesno_val,
         doit : no_op,
         set : checkerror_next_question }

     q3 := Question {
	 text : &account_parms_text[3],
	 response : &report_destination,
	 validate : reportdest_val,
         doit : no_op,
         set : reportdest_set }

     q4 := Question {
	 text : nil,
	 response : nil,
	 validate : no_op,
	 doit : prgexit,
	 set : no_op }

     account_parms.PushBack(q)
     account_parms.PushBack(q1)
     account_parms.PushBack(q2)
     account_parms.PushBack(q3)
     account_parms.PushBack(q4)
}

var account_errormess = make(map[int]errormesstype)

func init() {
     var err errormesstype
     msg := "Please enter S or R."
     err.message = &msg
     err.build = nil
     account_errormess[ENTER_S_OR_R] = err

     msg1 := "Please enter Y or N."
     err.message = &msg1
     account_errormess[ENTER_Y_OR_N] = err

     msg2 := "The starting account must be smaller than the ending account"
     err.message = &msg2
     account_errormess[START_ACCOUNT_LARGER] = err

     msg3 := "The account number must be between 100 and 1000"
     err.message = &msg3
     account_errormess[BAD_ACCOUNT_NUMBER] = err

     msg4 := "Please enter P, S, or D"
     err.message = &msg4
     account_errormess[ENTER_P_S_OR_D] = err

     msg5 := "That file already exists"
     err.message = &msg5
     account_errormess[FILE_EXISTS] = err

     msg6 := "Choose one of the following:"
     err.message = &msg6
     err.build = bld_bad_parmname
     account_errormess[BAD_PARM_NAME] = err
}

func main() {
     var errstat int
     var prctl *prcontrol = new(prcontrol)

     prctl.current_group = account_parms
     prctl.errormess = account_errormess

     if errstat = prompter(prctl); errstat > 0 {
	 handle_error(errstat, account_errormess)
     }

     // Print the report with gathered parameters
}

