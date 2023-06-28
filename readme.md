Tardigrade
==========

![Tardigrade](./img/tardig-11.png)

# Bookmark all your commands!!!

## Description

Tardigrade is a bookmark tool for commands. Tardigrade lets you easily navigate throughout dozens and dozens of commands. Tardigrade is a command line tool that assists remembering and retrieving any type of commands. One can store commands that had been used before. The commands are saved in a hierarchical and organized way that is easy to retrieve.

## Usage example

- normal usage, using tardigrade then finding and calling a command

![usage sample](https://s11.gifyu.com/images/SQas7.gif)

- sample yaml file:

```yaml
group1:
  group12:
  - echo some thing ^ echoing something
  - echo something else
  group14:
  - ls -la ^ list all files @ls @listing-files
  - echo hi there ^ echoing something
group2:
  group22:
  - pwd
  - echo com12 2 22 ^ @com2
  - echo com12 2 23 ^ @com2
  - echo com12 2 24 ^ @com2
  echo <> aa ^ this comm has variables @comm:
  - one # will convert to echo one aa
  - two # will convert to echo two aa
  echo <> and <> ^ this comm has several variables @comm:
  - a1, b1 # will convert to echo a1 and b1
  - a2, b2 # will convert to echo a2 and b2
```

- combo usage, command with several options

![usage sample](https://s11.gifyu.com/images/SQasj.gif)

- navigation, notice the footer description and options from groups

![usage sample](https://s11.gifyu.com/images/SQaxE.gif)

- going to another directory using tardigrade

![usage sample](https://s12.gifyu.com/images/SQaxg.gif)

- here is an older example of usage. Go here to see the example usage in a gif if github cannot display it: (https://s11.gifyu.com/images/SQ6mR.gif)

## More Description

If you have ever copied a bunch of commands in different notes or emails or readmes; and a month later or a year later or 3 years later; you need to look for that command; this tool will help you remember and reqtrieve them quickly. With Tardigrade all commands can be copied to a file in your user home directory, so all the commands that you have ever used, are in one place, and they can be easily retrieved from that file. But not stopping there, following what other tools do. The commands can be read from your user home directory and also from the current local directory where the Tardigrade app is called from and Tardigrade will combine all the keywords and commands from both directories. This concept can help any user, to personalize with commands to remember based on a directory. It can also help teams and projects, to remember their most used commands. Tardigrade works in most terminals.

## Tutorial

After installation. Type "tg" (or "tt" if the command was added to an rc file). Tardigrade will show all the main options in a menu. In each menu, the user can go up and down. If the user presses right or presses enter then the group or command will be chosen. If a command is chosen then the command will run (or print) and tardigrade will end. If the user presses left then the menu will go back to its parent.

Another feature is filtering, instead of going up and down one can start typing the a part of the command and the filter will show only the commands that match.

Another feature is Tardigrade can be started in flat mode and in tag mode to search keywords only in tags and in all mode to search keywords in either the command, description or tags.

### Tardicontent

Tardigrade uses a yaml file called tardicontent.yml. An example is earlier in the readme file.

The file follows a yaml hierachical structure. Any group ends in a semicolon :, any command is inside a yaml list item. The big exception is if a group that has a semicolon has the @comm tag, then it will become a command, and each of its list members will become a replacement for the command. Any comment that will be taken by tardigrade can be added after the ^ symbol, before the colon : if it is inside a group, at the end if it is a list item. In addition yaml comments can be added at the very end #, but those comments wont be taken by tardigrade. A tag starts with an at sign @, in the comment section. Any tag will be taken by tardigrade and are hierarchical so all children will inherit a tag. A group with a tag comm, will become a command as specified earlier.

### Targdisettings

A tardisettings file is created at first, with a group called settings. Here are some important attributes:
```yaml
settings:
    height: 11
    historysize: 10
    footerkeymaxsize: 12
```
```
settings (description):
    height: height of the command window
    historysize: how many commands can be save in the history
    footerkeymaxsize: the maximum size of each footer option
```

### Tardihistory

Tardigrade has its own history file called targdihistory.yml. Its a tardicontent yaml file with one group called history. Any command used with tardigrade will get copied to the history as a first member of the group. The history file will take care that there are no repeated commands.

## Thanks

This project relies in great libraries like kong, viper, and many others. But the main library that is relying on are the charm libraries, lipgloss, gum. Mainly gum (https://github.com/charmbracelet/gum), I copied all the filter section to customize it to tardigrade's needs.

## Installation

As of now, there are only two ways to install. I will work on the installation process in the next few months so it becomes easier to more users.

First, via go. Download the repo and run:
```
make build 
```
and the executable will be inside the build directory, then copy the path of the executable to the user path.

Second, install by using the deb package from the release section. Download the deb package: (https://github.com/sebastianxyzsss/tardigrade/releases)

to install: 
```
sudo dpkg -i tg.<version>.deb
```
then, add this script to .bashrc or zshrc:
```bash
tt() {
  echo ----------------------------------------------------
  _resultcomm=$(tg $@)
  echo $_resultcomm
  echo ----------------------------------------------------
  # print -S "$_resultcomm" # if available, saves in history
  eval ${_resultcomm}
}
```
** the script makes Tardigrade run the command instead of just printing the command so Tardigrade is more useful used this way.
