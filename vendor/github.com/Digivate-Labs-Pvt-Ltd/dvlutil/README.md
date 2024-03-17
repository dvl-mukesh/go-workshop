# dvlutil Library

This Golang library comprises a valuable assortment of functions designed to streamline and optimize various time-consuming tasks, making it an indispensable resource for enhancing efficiency within Golang applications.

# `ReadEnvVars` Function

The `ReadEnvVars` function is a Go utility that simplifies the process of reading environment variables and populating a configuration struct in your application. This can be especially useful when working with configurations that are typically set via environment variables in containerized or cloud-native applications.

## Use Case

In many Go applications, you need to configure settings or parameters, such as API keys, database connections, or feature toggles. These configurations often come from environment variables, and it can be cumbersome to manually read and parse these variables, especially if you have numerous configuration options.

The `ReadEnvVars` function is designed to address this challenge. It allows you to specify environment variable names for configuration fields and automatically populates a struct with the values from those environment variables. It also supports marking variables as required, making it easy to handle missing or invalid configurations gracefully.

## How to Use

1. **Import the Package**

   First, import the `ReadEnvVars` function in your Go code by adding the following line at the top of your Go source file:

   ```go
   import "github.com/Digivate-Labs-Pvt-Ltd/dvlutil"
   ```

2. **Define a Configuration Struct**

   Define a Go struct that represents your application's configuration. This struct should have fields corresponding to the configuration values you want to read from environment variables. You can add struct tags to specify the corresponding environment variable name and optional attributes like "required."

   ```go
   type AppConfig struct {
       APIKey  string    `env:"API_KEY,required"`
       Debug   bool      `env:"DEBUG"`
       Port    string    `env:"PORT"`
   }
   ```

   In this example, the `APIKey` field is marked as "required," and the `DEBUG` and `Port` fields are optional.

3. **Call `ReadEnvVars`**

   Call the `ReadEnvVars` function with a pointer to an instance of your configuration struct. This function will populate the struct fields with values from the specified environment variables.

   ```go
   cfg := &AppConfig{}
   err := ReadEnvVars(cfg)
   if err != nil {
       // Handle errors, e.g., missing required variables
   }
   ```

4. **Handling Errors**

   If any required environment variables are missing, the function will return an error. You should handle this error and take appropriate action, such as logging the error and exiting the application if necessary.

   ```go
   if err != nil {
       log.Fatalf("Error: %v", err)
   }
   ```

   If non-required variables are missing, the corresponding fields in your struct will be left with their zero values.

5. **Access Configuration Values**

   You can now access your configuration values directly from the struct fields. For example:

   ```go
   fmt.Printf("API Key: %s\n", cfg.APIKey)
   fmt.Printf("Debug Mode: %v\n", cfg.Debug)
   fmt.Printf("Port: %s\n", cfg.Port)
   ```

## Example

Here's an example of how to use the `ReadEnvVars` function with a sample configuration struct:

```go
package main

import (
    "fmt"
    "log"
    "github.com/Digivate-Labs-Pvt-Ltd/dvlutil"
)

type AppConfig struct {
    APIKey  string     `env:"API_KEY,required"`
    Debug   bool       `env:"DEBUG"`
    Port    string     `env:"PORT"`
}

func main() {
    cfg := &AppConfig{}
    err := dvlutil.ReadEnvVars(cfg)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("API Key: %s\n", cfg.APIKey)
    fmt.Printf("Debug Mode: %v\n", cfg.Debug)
    fmt.Printf("Port: %s\n", cfg.Port)
}
```

In this example, the `API_KEY` environment variable is marked as required, so it must be set for the application to run. The `DEBUG` and `PORT` environment variables are optional.

By using the `ReadEnvVars` function and the `env` struct tags, you can easily manage your application's configuration via environment variables, improving flexibility and maintainability.

Feel free to adapt and extend the provided sample code and configuration struct to suit your application's specific needs.
