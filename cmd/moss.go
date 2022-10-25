/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	language     string
	extension    string
	userId       string
	problem      string
	comment      string
	maxLimit     int64
	experimental bool
	numberResult int64
	moss         utils.Moss
)

func addFile(path string, info fs.DirEntry, err error) error {
	if err != nil {
		log.Error(err)
		return err
	}

	if info.IsDir() {
		log.WithField("dir", info.IsDir()).Debug(path)
		return nil
	}

	file, err := utils.ResolvePath(path)
	if err != nil {
		log.Error(err)
		return nil
	}

	if file.Problem != problem {
		log.Debugf("SKIP %s", path)
		return nil
	}

	moss.AddFile(path)
	log.Debugf("Find %s", path)

	return nil
}

// mossCmd represents the moss command
var mossCmd = &cobra.Command{
	Use:   "moss <CodePath>",
	Short: "A online code revies tool",
	Long: `The -l option specifies the source language of the tested programs.
    Moss supports many different languages; see the variable "languages" below for the
    full list.
    ["c", "cc", "java", "ml", "pascal", "ada", "lisp", "scheme", "haskell", "fortran", "ascii", "vhdl", "perl", "matlab", "python", "mips", "prolog", "spice", "vb", "csharp", "modula2", "a8086", "javascript", "plsql"]

    Example: Compare the lisp programs foo.lisp and bar.lisp:

    moss -l lisp -p problem ./Path/to/code/dir

    The -p option sets the problem name that will be reviews.

    Example:

    moss -l cc -p problem ./codePath

    The -m option sets the maximum number of times a given passage may appear
    before it is ignored.  A passage of code that appears in many programs
    is probably legitimate sharing and not the result of plagiarism.  With -m N,
    any passage appearing in more than N programs is treated as if it appeared in
    a base file (i.e., it is never reported).  Option -m can be used to control
    moss' sensitivity.  With -m 2, moss reports only passages that appear
    in exactly two programs.  If one expects many very similar solutions
    (e.g., the short first assignments typical of introductory programming
    courses) then using -m 3 or -m 4 is a good way to eliminate all but
    truly unusual matches between programs while still being able to detect
    3-way or 4-way plagiarism.  With -m 1000000 (or any very
    large number), moss reports all matches, no matter how often they appear.
    The -m setting is most useful for large assignments where one also a base file
    expected to hold all legitimately shared code.  The default for -m is 10.

    Examples:

    moss -l pascal -p problem -m 2 ./Path/to/code/dir
    moss -l cc -p problem -m 1000000 ./Path/to/code/dir

    The -c option supplies a comment string that is attached to the generated
    report.  This option facilitates matching queries submitted with replies
    received, especially when several queries are submitted at once.

    Example:

    moss -l scheme -p problem -c "Scheme programs" ./Path/to/code/dir

    The -n option determines the number of matching files to show in the results.
    The default is 250.

    Example:
    moss -l java -p problem -n 200 ./Path/to/code/dir

    The -x option sends queries to the current experimental version of the server.
    The experimental server has the most recent Moss features and is also usually
    less stable (read: may have more bugs).

    Example:

    moss -x -l ml -p problem ./Path/to/code/dir
    `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			codePath = viper.GetString("CodePath")
			if codePath == "" {
				err := errors.New("Requires args of Code Path")
				log.Error(err)
				return err
			}
			return nil
		}
		codePath = args[0]
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("ReviewProblem", cmd.Flags().Lookup("problem"))
		viper.BindPFlag("ReviewUserID", cmd.Flags().Lookup("userid"))
		viper.BindPFlag("ReviewLanguage", cmd.Flags().Lookup("language"))
		viper.BindPFlag("ReviewComment", cmd.Flags().Lookup("comment"))
		viper.BindPFlag("ReviewMaxLimit", cmd.Flags().Lookup("maxlimit"))
		viper.BindPFlag("ReviewExperimental", cmd.Flags().Lookup("experimental"))
		viper.BindPFlag("ReviewNumberResult", cmd.Flags().Lookup("numberresult"))

		problem = viper.GetString("ReviewProblem")
		userId = viper.GetString("ReviewUserID")
		language = viper.GetString("ReviewLanguage")
		comment = viper.GetString("ReviewComment")
		maxLimit = viper.GetInt64("ReviewMaxLimit")
		experimental = viper.GetBool("ReviewExperimental")
		numberResult = viper.GetInt64("ReviewNumberResult")

		log.Info("CodePath: ", codePath)
		log.Info("ReviewProblem: ", problem)
		log.Info("ReviewUserID: ", userId)
		log.Info("ReviewLanguage: ", language)
		log.Info("ReviewComment: ", comment)
		log.Info("ReviewMaxLimit: ", maxLimit)
		log.Info("ReviewExperimental: ", experimental)
		log.Info("ReviewNumberResult: ", numberResult)
	},
	Run: func(cmd *cobra.Command, args []string) {

		exp := 0

		if experimental {
			exp = 1
		}

		moss = utils.Moss{
			UserId:       userId,
			Language:     language,
			MaxLimit:     maxLimit,
			Comment:      comment,
			Experimental: exp,
			NumberResult: numberResult,
			Files:        []string{},
		}

		filepath.WalkDir(codePath, addFile)

		log.Infof("Find [%d] files, uploading...", len(moss.Files))

		result := moss.Review()

		log.Info(result)
	},
}

func init() {
	rootCmd.AddCommand(mossCmd)

	mossCmd.Flags().StringVarP(&problem, "problem", "p", "", "The problem to be reviewed")
	mossCmd.Flags().StringVarP(&userId, "userid", "u", "", "Review account")
	mossCmd.Flags().StringVarP(&language, "language", "l", "cc", "The reviewed code language")
	mossCmd.Flags().StringVarP(&comment, "comment", "c", "", "The comment")
	mossCmd.Flags().Int64VarP(&maxLimit, "maxlimit", "m", 10, "The maximum times before ignored")
	mossCmd.Flags().BoolVarP(&experimental, "experimental", "e", false, "Use experimental version of the server")
	mossCmd.Flags().Int64VarP(&numberResult, "numberresult", "n", 250, "The number of matching files to show in results")
}
