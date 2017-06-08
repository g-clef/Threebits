package structures

import (
    "strconv"
    )

type Response struct{
	Success bool;
	Test Test;
	Message string;
	AuthKey string;
}

type Test struct {
	Target string;
	Port int;
	Test string;
	Args Arguments;
}


type Arguments struct {
    Generic_TCP struct {
        Send string
        Receive string
        }
    Generic_UDP struct {
        Send string
        Receive string 
        }
}

func (t *Test) GetAddr() (string) {
	return t.Target + ":" + strconv.Itoa(t.Port)
}

func (r *Response) Stringify() (string){
	return r.Test.Target + "," + strconv.Itoa(r.Test.Port) + "," + r.Test.Test + "," + strconv.FormatBool(r.Success) + "," + r.Message
}


var AllArgs Arguments

