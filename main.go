package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

const PACK_RESULUTION int = 64

func main() {
	if err := CleanTargetPictures(); err != nil {
		log.Fatal(err)
	}

	packName := ""

	err := PromptUser("Enter Pack Name: ", &packName)
	if err != nil {
		log.Fatal(err)
	}

	packName, err = CreateTexturesPackFolder(packName)
	if err != nil {
		log.Fatal(err)
	}

	relaitvePath := "./" + packName

	err = CopyFile("./1.21.9-Template/pack.mcmeta", relaitvePath+"/"+"pack.mcmeta")
	if err != nil {
		log.Fatal(err)
	}

	err = CopyFile("./targetPictures/jamie.png", relaitvePath+"/"+"pack.png")
	if err != nil {
		log.Fatal(err)
	}

	relaitvePath, err = CreateDir(relaitvePath, "assets")
	if err != nil {
		log.Fatal(err)
	}

	relaitvePath, err = CreateDir(relaitvePath, "minecraft")
	if err != nil {
		log.Fatal(err)
	}

	relaitvePath, err = CreateDir(relaitvePath, "textures")
	if err != nil {
		log.Fatal(err)
	}

	relaitvePath, err = CreateDir(relaitvePath, "block")
	if err != nil {
		log.Fatal(err)
	}

	err = GenerateTargetImages("./1.21.9-Template/assets/minecraft/textures/block", relaitvePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done! :)")
}

func CleanTargetPictures() error {
	entries, err := os.ReadDir("./pictures")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		src, err := imaging.Open("./pictures/" + entry.Name())
		if err != nil {
			return err
		}

		src = imaging.Resize(src, PACK_RESULUTION, PACK_RESULUTION, imaging.Lanczos)
		src = imaging.AdjustSaturation(src, -50)

		err = imaging.Save(src, "./targetPictures/" + entry.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func PromptUser(prompt string, output *string) error {
	fmt.Printf("%s", prompt)
	_, err := fmt.Scanln(output)
	return err
}

func CreateDir(location, dirName string) (string, error) {
	location += "/" + dirName

	err := os.Mkdir(location, os.ModePerm)
	if err != nil {
		return "", err
	}

	return location, nil
}

func CreateTexturesPackFolder(targetPackName string) (string, error) {
	err := os.Mkdir(targetPackName, os.ModePerm)

	if err == nil {
		return targetPackName, nil
	}

	if !os.IsExist(err) {
		return "", err
	}

	fmt.Printf("Pack Name is already is use. ")

	targetPackName = targetPackName + "_" + strconv.Itoa(rand.Int())
	fmt.Printf("Pack will be renamed to: %s\n", targetPackName)

	userInput := ""

	if err := PromptUser("Proceed? (y/n): ", &userInput); err != nil {
		return "", err
	}

	userInput = strings.TrimSpace(userInput)

	if len(userInput) == 0 || strings.ToLower(userInput)[0] != 'y' {
		return "", errors.New("failed to create texture pack folder")
	}

	return CreateTexturesPackFolder(targetPackName)
}

func CopyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	// Ensure the source file is closed when the function returns
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	// Ensure the destination file is closed when the function returns
	defer destinationFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush the destination file to ensure all data is written to disk
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func GenerateTargetImages(textureNameLocation, dist string) error {
	targetPics, err := os.ReadDir("./targetPictures")
	if err != nil {
		return err
	}

	textureNames, err := os.ReadDir(textureNameLocation)
	if err != nil {
		return err
	}

	for _, textureName := range textureNames {
		if textureName.IsDir() {
			continue
		}

		name := textureName.Name()

		if !strings.HasSuffix(name, ".png") {
			continue
		}

		srcImg, err := imaging.Open(textureNameLocation + "/" + textureName.Name())
		if err != nil {
			return err
		}

		srcImg = imaging.Resize(srcImg, PACK_RESULUTION, PACK_RESULUTION, imaging.Lanczos)

		overlayImg, err := imaging.Open("./targetPictures/" + targetPics[rand.Intn(len(targetPics))].Name())
		if err != nil {
			return err
		}

		output := OverlayWithHoles(srcImg, overlayImg)
		output = imaging.OverlayCenter(srcImg, output, 0.5)

		err = imaging.Save(output,  dist + "/" + textureName.Name())

		if err != nil {
			return err
		}
	}

	return nil
}

func OverlayWithHoles(srcImg, overlayImg image.Image) *image.NRGBA {
	output := imaging.New(PACK_RESULUTION, PACK_RESULUTION, color.Transparent)

	for x := range PACK_RESULUTION {
		for y := range PACK_RESULUTION {
			srcCol := srcImg.At(x, y)

			if _, _, _, srcAlpha := srcCol.RGBA(); srcAlpha < 50 {
				continue
			}

			output.Set(x, y, overlayImg.At(x, y))
		}
	}

	return output
}