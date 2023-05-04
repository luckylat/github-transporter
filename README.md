# github-transporter

## Notice
- The product is not targeted to use GitHub/GitLab as storage.  
- It may drop working status, so it use carefully.
- It presuppose to exist remote git repository stably. if you use unstable remote server, don't use this.

## What is this?
- save git data only `.github-transporter`.

## Why I Created?
I using obsidian on OneDrive with plugins.  
some plugin has `.git` file, so the time to sync is so long.

## Functions
### Required
- Current folder is already git repository and git origin url is existed

### export
```
github-transport export
```
- make `.github-transporter` file.  
- remove `.git` folder.  

**local modified data will be lost, must push to origin.**

### import
```
github-transport import
```
- make `.git` folder.
- remove `.github-transporter` file (because it is unnecessary in this state).

