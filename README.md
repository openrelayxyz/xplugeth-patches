# xplugeth-patches

This repository contains a collection of patches used by Plugeth2 and the
xplugeth build tool to support patching various Geth forks with the PluGeth2
plugin hooks.

Unless you know why you're here, this project probably isn't what you're looking
for. You'll find far more useful information at the 
[xplugeth](https://github.com/openrelayxyz/xplugeth) git repository.

If you're looking to make your own extensions to PluGeth2, here's how this repository
works:

Every branch on this repository starts with a commit from some Geth fork, initially
ethereum/go-ethereum, maticnetwork/bor, or etclabscore/core-geth. We then make a small
change to incorporate xplugeth hooks into that change, and commit it in a single Git
commit. That git commit can then be cherry-picked onto Geth forks by the xplugeth build
tool to enable those hooks, essentially applying a simple patch.

The branches / tags here are essentially considered read-only. Each branch should only
contain one commit beyond its upstream commit. If future changes are needed, we handle 
that by creating a new branch with a new single commit. These patches are designed
to be fairly resilient to changes in upstream code. We do not commit updates to the
go.mod or go.sum as those would drastically increase the likelihood of conflicts in
applying the patch to later versions of Geth (instead, the go.mod and go.sum are
updated at build time by the xplugeth build tool).
