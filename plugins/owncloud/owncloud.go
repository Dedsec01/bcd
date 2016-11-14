package owncloud

import (
        log "github.com/Sirupsen/logrus"
        "github.com/fsouza/go-dockerclient"
        "github.com/bytesizedhosting/bcd/plugins"
        "net/rpc"
)

type Owncloud struct {
        plugins.Base
        imageName string
}

func New(client *docker.Client) (*Owncloud, error) {
        manifest, err := plugins.LoadManifest("owncloud")

        if err != nil {
                return nil, err
        }

        return &Owncloud{Base: plugins.Base{DockerClient: client, Name: "owncloud", Version: 1, Manifest: manifest}, imageName: "bytesized/owncloud"}, nil
}

func (self *Owncloud) RegisterRPC(server *rpc.Server) {
        rpc := plugins.NewBaseRPC(self)
        server.Register(&OwncloudRPC{base: self, BaseRPC: *rpc})
}

type OwncloudOpts struct {
        plugins.BaseOpts
}

func (self *Owncloud) Install(opts *OwncloudOpts) error {
        var err error

        err = opts.SetDefault(self.Name)
        if err != nil {
                return err
        }

        log.WithFields(log.Fields{
                "plugin":       self.Name,
                "config_folder": opts.ConfigFolder,
                "username":     opts.Username,
        }).Debug("Plugin options")

        log.Debugln("Pulling docker image", self.imageName)
        err = self.DockerClient.PullImage(docker.PullImageOptions{Repository: self.imageName}, docker.AuthConfiguration{})
        if err != nil {
                return err
        }

        portBindings := map[docker.Port][]docker.PortBinding{
                "8989/tcp": []docker.PortBinding{docker.PortBinding{HostPort: opts.WebPort}},
        }

        hostConfig := docker.HostConfig{
                PortBindings: portBindings,
                Binds:        plugins.DefaultBindings(opts),
        }

        conf := docker.Config{Env: []string{"PUID=" + opts.User.Uid}, Image: self.imageName}

        log.Debugln("Creating docker container")
        c, err := self.DockerClient.CreateContainer(docker.CreateContainerOptions{Config: &conf, HostConfig: &hostConfig, Name: "bytesized_owncloud_" + opts.WebPort})

        if err != nil {
                return err
        }

        log.Debugln("Starting docker container")

        err = self.DockerClient.StartContainer(c.ID, nil)
        if err != nil {
                return err
        }

        opts.ContainerId = c.ID

        return nil
}
