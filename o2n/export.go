package o2n

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
	utilerrors "k8s.io/kubernetes/pkg/util/errors"

	cclientcmd "github.com/openshift/origin/pkg/cmd/util/clientcmd"
)

func Export(outputDir string, f *cclientcmd.Factory, cmd *cobra.Command, args []string) error {
	jww.DEBUG.Printf("Starting Export")
	jww.DEBUG.Printf("outputDir = %s", outputDir)
	jww.DEBUG.Printf("args = %s", args)

	// based on github.com/openshift/origin/pkg/cmd/cli/cmd/export.go

	outputVersion := "v1"

	exporter := &myExporter{}

	cmdNamespace, _, err := f.DefaultNamespace()
	if err != nil {
		return err
	}

	mapper, typer := f.Object()
	b := resource.NewBuilder(mapper, typer, f.ClientMapperForCommand()).
		NamespaceParam(cmdNamespace).
		//.DefaultNamespace().AllNamespaces(allNamespaces).
		//FilenameParam(explicit, filenames...).
		ResourceTypeOrNameArgs(true, args...).
		Flatten()

	one := false
	infos, err := b.Do().IntoSingular(&one).Infos()
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		return fmt.Errorf("no resources found - nothing to export")
	}

	// remove runtime information from resources
	newInfos := []*resource.Info{}
	errs := []error{}
	for _, info := range infos {
		if err := exporter.Export(info.Object, false); err != nil {
			if err == ErrExportOmit {
				continue
			}
			errs = append(errs, err)
		}
		newInfos = append(newInfos, info)
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}
	infos = newInfos

	var result runtime.Object
	object, err := resource.AsVersionedObject(infos, !one, outputVersion)
	if err != nil {
		return err
	}

	result = object

	p, _, err := kubectl.GetPrinter("json", "")
	if err != nil {
		return err
	}
	out := os.Stdout
	return p.PrintObj(result, out)

}
