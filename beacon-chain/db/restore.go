package db

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/kv"
	"github.com/prysmaticlabs/prysm/shared/cmd"
	"github.com/prysmaticlabs/prysm/shared/fileutil"
	"github.com/prysmaticlabs/prysm/shared/promptutil"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const dbExistsYesNoPrompt = "A database file already exists in the target directory. " +
	"Are you sure that you want to overwrite it? [y/n]"

func restore(cliCtx *cli.Context) error {
	sourceFile := cliCtx.String(cmd.RestoreSourceFileFlag.Name)
	targetDir := cliCtx.String(cmd.RestoreTargetDirFlag.Name)

	restoreDir := path.Join(targetDir, kv.BeaconNodeDbDirName)
	if fileutil.FileExists(path.Join(restoreDir, kv.DatabaseFileName)) {
		resp, err := promptutil.ValidatePrompt(
			os.Stdin, dbExistsYesNoPrompt, promptutil.ValidateYesOrNo,
		)
		if err != nil {
			return errors.Wrap(err, "could not validate choice")
		}
		if strings.ToLower(resp) == "n" {
			logrus.Info("Restore aborted")
			return nil
		}
	}
	if err := fileutil.MkdirAll(restoreDir); err != nil {
		return err
	}
	if err := fileutil.CopyFile(sourceFile, path.Join(restoreDir, kv.DatabaseFileName)); err != nil {
		return err
	}

	logrus.Info("Restore completed successfully")
	return nil
}
