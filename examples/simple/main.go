package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yourneighborhoodchef/browserforge/fingerprint"
)

func main() {

	generator, err := fingerprint.New()
	if err != nil {
		log.Fatalf("Error initializing fingerprint generator: %v", err)
	}

	fp, err := generator.Generate()
	if err != nil {
		log.Fatalf("Error generating fingerprint: %v", err)
	}

	fpJSON, _ := json.MarshalIndent(fp, "", "  ")
	fmt.Println("Generated Fingerprint:")
	fmt.Println(string(fpJSON))

	headers, err := generator.GenerateHeadersOnly()
	if err != nil {
		log.Fatalf("Error generating headers: %v", err)
	}

	fmt.Println("\nGenerated Headers Only:")
	headersJSON, _ := json.MarshalIndent(headers, "", "  ")
	fmt.Println(string(headersJSON))

	customGenerator, err := fingerprint.NewWithOptions(
		fingerprint.WithBrowser("chrome"),
		fingerprint.WithOperatingSystem("windows"),
	)
	if err != nil {
		log.Fatalf("Error initializing custom generator: %v", err)
	}

	customFp, err := customGenerator.Generate()
	if err != nil {
		log.Fatalf("Error generating custom fingerprint: %v", err)
	}

	fmt.Println("\nCustom Fingerprint Details:")
	fmt.Printf("User Agent: %s\n", customFp.Navigator.UserAgent)
	fmt.Printf("Browser: %s\n", customFp.Navigator.AppVersion)
	if customFp.Navigator.Oscpu != nil {
		fmt.Printf("Operating System: %s\n", *customFp.Navigator.Oscpu)
	}

	fmt.Println("\nFew HTTP Headers:")
	for key, value := range customFp.Headers {
		if key == "User-Agent" || key == "Accept-Language" || key == "Accept" {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	camoufoxGenerator, err := fingerprint.NewWithOptions(
		fingerprint.WithCamoufoxConstraints(),
	)
	if err != nil {
		log.Fatalf("Error initializing Camoufox generator: %v", err)
	}

	camoufoxGenerator.SetFirefoxVersion("115.0")

	camoufoxFp, err := camoufoxGenerator.Generate()
	if err != nil {
		log.Fatalf("Error generating Camoufox fingerprint: %v", err)
	}

	fmt.Println("\nCamoufox Fingerprint Details:")
	fmt.Printf("User Agent: %s\n", camoufoxFp.Navigator.UserAgent)
	fmt.Printf("Browser: %s\n", camoufoxFp.Navigator.AppVersion)
	if camoufoxFp.Navigator.Oscpu != nil {
		fmt.Printf("Operating System: %s\n", *camoufoxFp.Navigator.Oscpu)
	}

	fmt.Println("\nCamoufox Screen Details:")
	fmt.Printf("Screen: %dx%d\n", camoufoxFp.Screen.Width, camoufoxFp.Screen.Height)
	fmt.Printf("Available: %dx%d\n", camoufoxFp.Screen.AvailWidth, camoufoxFp.Screen.AvailHeight)
	fmt.Printf("Window Outer: %dx%d\n", camoufoxFp.Screen.OuterWidth, camoufoxFp.Screen.OuterHeight)
	fmt.Printf("Window Inner: %dx%d\n", camoufoxFp.Screen.InnerWidth, camoufoxFp.Screen.InnerHeight)
	fmt.Printf("Window Position: (%d, %d)\n", camoufoxFp.Screen.ScreenX, camoufoxFp.Screen.PageYOffset)

	constrainedGenerator, err := fingerprint.NewWithOptions(
		fingerprint.WithCamoufoxConstraints(),
		fingerprint.WithScreenConstraints(1920, 1080),
		fingerprint.WithWindowSize(1200, 800),
	)
	if err != nil {
		log.Fatalf("Error initializing constrained generator: %v", err)
	}

	constrainedFp, err := constrainedGenerator.Generate()
	if err != nil {
		log.Fatalf("Error generating constrained fingerprint: %v", err)
	}

	fmt.Println("\nConstrained Screen Details:")
	fmt.Printf("Screen: %dx%d\n", constrainedFp.Screen.Width, constrainedFp.Screen.Height)
	fmt.Printf("Available: %dx%d\n", constrainedFp.Screen.AvailWidth, constrainedFp.Screen.AvailHeight)
	fmt.Printf("Window Outer: %dx%d\n", constrainedFp.Screen.OuterWidth, constrainedFp.Screen.OuterHeight)
	fmt.Printf("Window Inner: %dx%d\n", constrainedFp.Screen.InnerWidth, constrainedFp.Screen.InnerHeight)
	fmt.Printf("Window Position: (%d, %d)\n", constrainedFp.Screen.ScreenX, constrainedFp.Screen.PageYOffset)
}
