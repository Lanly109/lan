/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lanly109/lan/utils"
	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	srcPath       string
	srcPathLen    int
	problemList   []string
	problemSet    mapset.Set
	extensionList []string
	extensionSet  mapset.Set
	copyNum       int
)

func copyFile(src string, info fs.DirEntry, err error) error {
	// test/GD-00510/expr/expr.cpp
	if err != nil {
		log.Error(err)
		return err
	}

	if src == srcPath {
		return nil
	}

	filePath := src[srcPathLen+1:]

	if info.IsDir() {
		tmp := strings.Split(filePath, string(os.PathSeparator))
		if len(tmp) > 2 {
			log.Warnf("Invalid len of dir path: [%s]. Will skip it", src)
			return fs.SkipDir
		} else if len(tmp) == 2 { // GD-xxxx/problem
			if !problemSet.Contains(tmp[1]) { // tmp[1] -> problem name
				log.Warnf("Invalid problem[%s] of dir path: [%s]", tmp[1], src)
				return nil
			}
		}

		dest := filepath.Join(codePath, filePath)

		err = utils.Copy(src, dest)

		if err != nil {
			log.Errorf("Copy dir  [%s] -> [%s], error: %s", src, dest, err)
			return nil
		}

		log.Debugf("Copy dir  [%s] -> [%s]", src, dest)

		return nil
	}

	if filepath.Ext(filePath) == ".txt" {
		log.Debugf("Skip [%s]", src)
		return nil
	}

	codeInfo, err := utils.ResolvePath(filePath)
	if err != nil {
		log.Error(err)
		return nil
	}

	if !problemSet.Contains(codeInfo.Problem) {
		log.Debugf("Invalid problem[%s] of code file path: [%s]", codeInfo.Problem, src)
		return nil
	}

	if !extensionSet.Contains(codeInfo.Extension) {
		log.Debugf("Invalid extension[%s] of code file path: [%s]", codeInfo.Extension, src)
		return nil
	}

	dest := filepath.Join(codePath, filePath)

	err = utils.Copy(src, dest)

	if err != nil {
		log.Errorf("Copy file [%s] -> [%s], error: %s", src, dest, err)
		return nil
	}

	log.Infof("Copy file [%s] -> [%s]", src, dest)
	copyNum += 1
	return nil
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean <CodePath> <SourcePath>",
	Short: "Clear unnecessary files(except .cpp)",
	Long: `Clean unnecessary files, such as .exe .in .out .ans .pdf.
    Generate CCF form folder, only including .cpp.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {

			if len(args) < 1 {
				codePath = viper.GetString("CodePath")
				if codePath == "" {
					err := errors.New("Requires args of Code Path")
					log.Error(err)
					return err
				}
			} else {
				codePath = args[0]
			}

			srcPath = viper.GetString("SourcePath")

			if srcPath == "" {
				err := errors.New("Requires args of Source Path")
				log.Error(err)
				return err
			}

			return nil
		}
		codePath = args[0]
		srcPath = args[1]
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("Problems", cmd.Flags().Lookup("problems"))
		viper.BindPFlag("Extensions", cmd.Flags().Lookup("extensions"))

		problemList = viper.GetStringSlice("Problems")
		extensionList = viper.GetStringSlice("Extensions")

		srcPath = filepath.Clean(srcPath)
		srcPathLen = len(srcPath)
		codePath = filepath.Clean(codePath)

		log.Info("SourcePath: ", srcPath)
		log.Info("CodePath: ", codePath)
		log.Info("Problems: ", problemList)
		log.Info("Extentions: ", extensionList)

		var tmpList1, tmpList2 []interface{}
		for _, data := range problemList {
			tmpList1 = append(tmpList1, data)
		}
		for _, data := range extensionList {
			tmpList2 = append(tmpList2, data)
		}

		problemSet = mapset.NewSetFromSlice(tmpList1)
		extensionSet = mapset.NewSetFromSlice(tmpList2)
	},
	Run: func(cmd *cobra.Command, args []string) {
		copyNum = 0

		filepath.WalkDir(srcPath, copyFile)

		log.Infof("Copy %d files", copyNum)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().StringSliceVarP(&problemList, "problems", "", []string{"problem1", "problem2"}, "competition problems")
	cleanCmd.Flags().StringSliceVarP(&extensionList, "extensions", "", []string{".cpp", ".c", ".pas"}, "accepted code extensions")
}
