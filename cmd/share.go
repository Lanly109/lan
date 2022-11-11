/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var scriptName string

const sharingCode string = `#include <windows.h>
#include <cstdlib>
#include <time.h>
#include <cstdio>
#include <winuser.h>

void down(int vk){
	keybd_event(vk, 0, 0, 0);
}

void up(int vk){
	keybd_event(vk, 0, KEYEVENTF_KEYUP, 0);
}

void press(int vk){
	down(vk);
	Sleep(1);
	up(vk);
}

int main(){
	system("explorer D:");	
	Sleep(1000);
	down(VK_LWIN);
	press(VK_UP);
	up(VK_LWIN);
	system("explorer \\\\10.78.30.99\\b404"); // <please change the address>
	Sleep(1000);
	down(VK_LWIN);
	press(VK_RIGHT);
	up(VK_LWIN);
}
`

// shareCmd represents the sharing command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Generate script that is easy to use shared folders to collect code",
	Long:  `A cpp source code that will open the sharing directory as well as D:\ and Split screen between left and right to easy collect code.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("filename: %s", scriptName)
		if utils.FileExist(scriptName) {
			log.Errorf("The File %s Exists!", scriptName)
			return
		}
		err := utils.WriteFile(scriptName, sharingCode)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Successfully generated sharing script [%s]", scriptName)
		}
	},
}

func init() {
	genCmd.AddCommand(shareCmd)

	shareCmd.Flags().StringVarP(&scriptName, "name", "n", "sharing.cpp", "The script name(default sharing.cpp)")
}
