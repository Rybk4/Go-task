package templates

import "io"

// TemplateExecutor defines the interface for executing templates.
type TemplateExecutor interface {
	// ExecTemplate executes the template with the given name and data, writing the result to writer.
	ExecTemplate(writer io.Writer, name string, data interface{}) (err error)

	// ExecTemplateWithFunc executes the template with the given name and data, using the provided handler function.
	ExecTemplateWithFunc(writer io.Writer, name string, data interface{}, handlerFunc InvokeHandlerFunc) (err error)
}

// InvokeHandlerFunc defines the type of a handler function that can be called within a template.
type InvokeHandlerFunc func(handlerName string, methodName string, args ...interface{}) interface{}
