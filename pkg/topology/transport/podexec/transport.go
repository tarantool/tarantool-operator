package podexec

import (
	"bytes"
	"context"
	"time"

	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec/cli"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	resource       = "pods"
	subResource    = "exec"
	defaultTimeout = 2 * time.Second
)

type PodExec struct {
	RestClient    rest.Interface
	RestConfig    *rest.Config
	RuntimeScheme *runtime.Scheme

	CLI cli.CLI

	ContainerName string
}

func (r *PodExec) Exec(ctx context.Context, pod *v1.Pod, res interface{}, lua string, args ...interface{}) error {
	command, err := r.CLI.CreateCommand(lua, args...)
	if err != nil {
		return err
	}

	stdout, err := r.execShellCommand(ctx, command, pod.GetName(), pod.GetNamespace())
	if err != nil {
		return err
	}

	return r.CLI.Unmarshal(stdout, &res)
}

func (r *PodExec) execShellCommand(
	ctx context.Context,
	command *cli.Command,
	podName string,
	podNamespace string,
) (string, error) {
	stdinExist := false

	if command.StdIn != "" {
		stdinExist = true
	}

	execCreateLua := r.RestClient.
		Post().
		Namespace(podNamespace).
		Resource(resource).
		Name(podName).
		SubResource(subResource).
		Timeout(defaultTimeout).
		MaxRetries(1).
		VersionedParams(
			&v1.PodExecOptions{
				Stdout:    true,
				Stdin:     stdinExist,
				Stderr:    false,
				TTY:       false,
				Command:   command.Command,
				Container: r.ContainerName,
			},
			runtime.NewParameterCodec(r.RuntimeScheme),
		)

	exec, err := remotecommand.NewSPDYExecutor(r.RestConfig, "POST", execCreateLua.URL())
	if err != nil {
		return "", err
	}

	stdout := bytes.NewBufferString("")

	if stdinExist {
		stdin := bytes.NewBufferString(command.StdIn)
		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: stdout,
			Stdin:  stdin,
			Stderr: nil,
			Tty:    false,
		})
	} else {
		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: stdout,
			Stdin:  nil,
			Stderr: nil,
			Tty:    false,
		})
	}

	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}
