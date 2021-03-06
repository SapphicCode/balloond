# balloond

An automatic libvirt memory balloon daemon.

## Building

```
go get github.com/Pandentia/balloond
go build github.com/Pandentia/balloond/cmd/balloond
```

## Usage

A really basic example:

`./balloond -unix /var/run/libvirt-sock`

That's it. No, seriously.

## Explanation

**So, what's it do?**

To answer that question, you need to understand libvirt memory ballooning.

In a nutshell, libvirt memory ballooning *exists* but is an essentially manual process ("virsh setmem"). That's where this daemon comes in. It keeps a record of all running domains on a given system, and dynamically changes memory allocation every 10 seconds (configurable) where necessary. This ensures that all VMs (with a balloon driver) only use the memory they're actually using, while allowing the rest to be reaped by the hypervisor. It also gives each VM an amount of memory that will be guaranteed free and available, up to the maximum allowed memory by the domain's config.

## Future plans

See [#1](https://github.com/Pandentia/balloond/issues/1).

## Stability

This project is in its **alpha** state. Feel free to use at your own risk, but if all of your VMs suddenly go out of memory due to some unknown bug, don't blame me.

That said, if you do end up using it on your homelab and encounter a bug, please do feel encouraged to [report it](https://github.com/Pandentia/balloond/issues/new).
