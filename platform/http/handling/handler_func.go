package handling

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"platform/http/actionresults"
	"platform/services"
	"platform/templates"
	"reflect"
	"strings"
)

func createInvokehandlerFunc(ctx context.Context, routes []Route) templates.InvokeHandlerFunc {
	return func(handlerName, methodName string, args ...interface{}) interface{} {
		var err error

		// Loop through registered routes to find matching handler and method
		for _, route := range routes {
			if strings.EqualFold(handlerName, route.handlerName) &&
				strings.EqualFold(methodName, route.handlerMethod.Name) {

				// Prepare parameters for handler method
				paramVals := make([]reflect.Value, len(args))
				for i := 0; i < len(args); i++ {
					paramVals[i] = reflect.ValueOf(args[i])
				}

				// Create a new instance of the handler struct
				structVal := reflect.New(route.handlerMethod.Type.In(0))
				services.PopulateForContext(ctx, structVal.Interface())

				// Append the struct instance to the beginning of paramVals
				paramVals = append([]reflect.Value{structVal.Elem()}, paramVals...)

				// Call the handler method with the prepared parameters
				result := route.handlerMethod.Func.Call(paramVals)

				// Check if the result is a TemplateActionResult
				if action, ok := result[0].Interface().(*actionresults.TemplateActionResult); ok {
					// Recursive invocation to create a new invoke handler function
					invoker := createInvokehandlerFunc(ctx, routes)

					// Populate context for the action, including the invoker as an extra
					err = services.PopulateForContextWithExtras(ctx, action, map[reflect.Type]reflect.Value{
						reflect.TypeOf(invoker): reflect.ValueOf(invoker),
					})

					// Create a stringResponseWriter to capture template execution result
					writer := &stringResponseWriter{Builder: &strings.Builder{}}
					if err == nil {
						err = action.Execute(&actionresults.ActionContext{
							Context:        ctx,
							ResponseWriter: writer,
						})
						if err == nil {
							return (template.HTML)(writer.Builder.String())
						}
					}
				} else {
					return fmt.Sprint(result[0])
				}
			}
		}

		if err == nil {
			err = fmt.Errorf("No route found for %v %v", handlerName, methodName)
		}
		panic(err)
	}
}

type stringResponseWriter struct {
	*strings.Builder
}

func (sw *stringResponseWriter) Write(data []byte) (int, error) {
	return sw.Builder.Write(data)
}

func (sw *stringResponseWriter) WriteHeader(statusCode int) {}

func (sw *stringResponseWriter) Header() http.Header {
	return http.Header{}
}
