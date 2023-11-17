package main

import (
  "strings"
)

//represents all nodes even though not all members are used
//for each type of node
//Go does not have classes and we need polymorphism
type Node struct {
  tag string  //the tag of the html element parsed
  text string //if it is a terminal, this is the text between the opening and closing tags
  parent *Node //the node/element that contains this node/element
  children []*Node //if it is a container, these are the nodes/elements it contains
  emitHTML func( *Node ) string //emit is synonymous with “print” or “convert to string” here
}


func createNode( tag string, parent *Node ) *Node {
  var result = &Node {
    tag: tag,
    value: "",
    parent: parent
  }

  switch tag {
  case "section":
    result.emitHTML = emitSection
    break
  case "div":
    result.emitHTML = emitDiv
    break
  case "ul":
    result.emitHTML = emitUl
    break
  case "body":
    result.emitHTML = emitBody
    break
  case "p":
    result.emitHTML = emitPara
    break
  case "li":
    result.emitHTML = emitLi
    break
  case "h1":
    result.emitHTML = emitHn
    break
  case "h2":
    result.emitHTML = emitHn
    break
  default:
    result.emitHTML = emitBody
}

 if ( parent != nil ) {
   parent.addChild( result )
 }
 return result

}

//this notation allows us to use a dot notation to call the function addChild
//ex. n.addChild( c )
func ( n *Node ) addChild( c *Node ) {
  n.children = append(n.children, c)
}

//Emitters to be assigned at creation
func emitSection( c *Node ) string {
  var str string = "<section class='section'>"
  for _, value := range c.children {
    str += value.emitHTML( value )
  }
  str += "</section>"
  return str
}
func emitDiv( c *Node ) string {
  var str string = "<div class='container'><button class='ignore' onClick= 'this.parentNode.appendChild(CreateParagraph())' }>Add Paragraph</button><button class='ignore' onClick= 'this.parentNode.appendChild(CreateList())' }>Add List</button>"
  for _, value := range c.children {
    str += value.emitHTML( value )
  }
  str += "</div>"
  return str
}
func emitUl( c *Node ) string {
  var str string = "<ul class='list'><button class='ignore' onClick= 'this.parentNode.appendChild(CreateItem())' }>Add Item</button>"
  for _, value := range c.children {
    str += value.emitHTML( value )
  }
  str += "</ul>"
  return str
}
func emitBody( c *Node ) string {
  var str string = ""
  for _, value := range c.children {
    str += value.emitHTML( value )
  }
  return str
}

func emitPara( t *Node ) string {
  return "<textarea class='paragraph'>" + t.text + "</textarea>"
}
func emitHn( t *Node ) string {
  return "<input type='text' class='heading' value='" + t.text + "'></input>"
}
func emitLi( t *Node ) string {
  return "<input type='text' class='item' value='" + t.text + "'></input>"
}

func wrapElement( emitted string ) string {
  return emitted
}


//detect an opening tag
// read all content until we see the same closing tag
// create the appropriate node
//    If terminal, we just feed in the read text
//    If container, we repeat the above recursively
// go until we reach the end of the string

//returns the content between an opening and closing tag
func getContent( tag string, str string ) string {
  end := strings.Index( str, "</" + tag + ">")

  if ( end == -1 ) {
    return str
  }

  return str[0:end]
}

//returns the text between <>
//Because all tags should be right next to one another
//we can assume all str just start with a tag
func readTag( str string ) string {
  // var buff string = ""
  // var i int = 1
  //
  // for i < len(str) && str[i] != '<' {
  //   buff += string(str[i])
  //   i++
  // }
  var end int = strings.Index( str, ">" )
  if ( end == -1 ) {
    return ""
  }
  return str[1 : end]
}


// Create Intermediate Representation
// Generate a Node parent that contains the parsed text as Nodes
func createIR( str string ) *Node {

  //prentend like we have a body
  // createNode( "body", nil )

  //reading a tag, reading content
  // while a tag is found
      //two cases
        //terminal node
            //create the node and store the content
        //container node
          //create the node and store the children through a recursive call

    //ending condition
      //there isn't more string to be read (the current position is at the end)
  return nil
}


// Given a string that has read only html tags, transform it into a string that has input tags
// Use all of the previously defined functions to achieve this
func TransformStr( str string ) string {
  return str
}
