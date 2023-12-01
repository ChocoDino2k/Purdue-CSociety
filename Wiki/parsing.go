package main

import (
  "strings"
)

//represents all nodes even though not all members are used
//for each type of node
//Go does not have classes and we need polymorphism
type Node struct {
  tag string
  value string
  parent *Node
  children []*Node
  emitHTML func( *Node ) string
}

func createNode( tag string, parent *Node ) *Node {

  var result = &Node {
    tag: tag,
    value: "",
    parent: parent }

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
  return "<textarea class='paragraph'>" + t.value + "</textarea>"
}
func emitHn( t *Node ) string {
  return "<input type='text' class='heading' value='" + t.value + "'></input>"
}
func emitLi( t *Node ) string {
  return "<input type='text' class='item' value='" + t.value + "'></input>"
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

func createIR( str string ) *Node {
  var curContainer *Node = createNode("body", nil)
  var remainingStr string = str
  var tag string = readTag( remainingStr )

  for len(tag) > 0 {
    remainingStr = remainingStr[ (len(tag) + 2): ]
    var buff string = getContent( tag, remainingStr )
    //code based on tag
    //terminals just get the content
    if (
      tag == "li" ||
      tag == "p" ||
      tag == "h1" ||
      tag == "h2" ) {
      n := createNode( tag, curContainer )
      n.value = buff
    } else if (
      tag == "section" ||
      tag == "div" ||
      tag == "ul" ) {
      n := createNode( tag, curContainer )
      n.children = createIR( buff ).children
    }

    //end
    // content + < + tag + />
    if ( len(buff) + len(tag) + 3 >= len(remainingStr) ) {
      remainingStr = ""
    } else {
      remainingStr = remainingStr[ (len(tag) + len(buff) + 3): ]
    }
    tag = readTag( remainingStr )
  }
  return curContainer

}

func TransformStr( str string ) string {
  n := createIR( str )
  var result string = ""

  for i, value := range n.children {
    if ( i == 0 ) { continue }
    result += value.emitHTML( value )
  }

  return result

}
