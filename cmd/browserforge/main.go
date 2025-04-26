package main

import (
   "encoding/json"
   "fmt"
   "os"

   "github.com/yourneighborhoodchef/browserforge/fingerprint"
)

func main() {
   if len(os.Args) < 2 {
       fmt.Fprintf(os.Stderr, "Usage: %s [headers|fingerprint|all]\n", os.Args[0])
       os.Exit(1)
   }
   cmd := os.Args[1]
   
   // Create the fingerprint generator
   generator, err := fingerprint.New()
   if err != nil {
       fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
       os.Exit(1)
   }
   
   switch cmd {
   case "headers":
       // Generate only HTTP headers
       hdrs, err := generator.GenerateHeadersOnly()
       if err != nil {
           fmt.Fprintf(os.Stderr, "Error generating headers: %v\n", err)
           os.Exit(1)
       }
       out, _ := json.MarshalIndent(hdrs, "", "  ")
       fmt.Println(string(out))

   case "fingerprint", "all":
       // Generate complete fingerprint
       fp, err := generator.Generate()
       if err != nil {
           fmt.Fprintf(os.Stderr, "Error generating fingerprint: %v\n", err)
           os.Exit(1)
       }
       
       if cmd == "all" {
           // For the 'all' command, output the complete fingerprint including headers
           out, _ := json.MarshalIndent(fp, "", "  ")
           fmt.Println(string(out))
       } else {
           // For 'fingerprint', omit the headers to maintain backward compatibility
           // Create a copy of the struct without the Headers field
           type fingerprintWithoutHeaders struct {
               UserAgent            string `json:"userAgent"`
               UserAgentData        map[string]interface{} `json:"userAgentData"`
               AppVersion           string `json:"appVersion"`
               OSCpu                *string `json:"oscpu"`
               Product              string `json:"product"`
               HardwareConcurrency  int `json:"hardwareConcurrency"`
               DeviceMemory         *int `json:"deviceMemory"`
               ExtraProperties      map[string]interface{} `json:"extraProperties"`
               Screen               fingerprint.ScreenFingerprint `json:"screen"`
               AudioCodecs          map[string]string `json:"audioCodecs"`
               VideoCodecs          map[string]string `json:"videoCodecs"`
               PluginsData          map[string]interface{} `json:"pluginsData"`
               MultimediaDevices    []string `json:"multimediaDevices"`
               Battery              map[string]interface{} `json:"battery"`
               Fonts                []string `json:"fonts"`
           }
           
           noHeaders := fingerprintWithoutHeaders{
               UserAgent:           fp.Navigator.UserAgent,
               UserAgentData:       fp.Navigator.UserAgentData,
               AppVersion:          fp.Navigator.AppVersion,
               OSCpu:               fp.Navigator.Oscpu,
               Product:             fp.Navigator.Product,
               HardwareConcurrency: fp.Navigator.HardwareConcurrency,
               DeviceMemory:        fp.Navigator.DeviceMemory,
               ExtraProperties:     fp.Navigator.ExtraProperties,
               Screen:              fp.Screen,
               AudioCodecs:         fp.AudioCodecs,
               VideoCodecs:         fp.VideoCodecs,
               PluginsData:         fp.PluginsData,
               MultimediaDevices:   fp.MultimediaDevices,
               Battery:             fp.Battery,
               Fonts:               fp.Fonts,
           }
           
           out, _ := json.MarshalIndent(noHeaders, "", "  ")
           fmt.Println(string(out))
       }

   default:
       fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
       os.Exit(1)
   }
}