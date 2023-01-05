package gpujob

import (
	"fmt"
	"mini-kubernetes/tools/def"
)

func slurmGenerator(config def.GPUSlurmConfig) string {
	content := ``
	content += "#!/bin/bash\n\n"
	content += fmt.Sprintf("#SBATCH --job-name=%s\n", config.JobName)
	content += fmt.Sprintf("#SBATCH --partition=%s\n", config.Partition)
	content += fmt.Sprintf("#SBATCH -N %d\n", config.Node)
	content += fmt.Sprintf("#SBATCH --cpus-per-task=%d\n", config.CpusPerTask)
	content += fmt.Sprintf("#SBATCH --ntasks-per-node=%d\n", config.NtasksPerNode)
	content += fmt.Sprintf("#SBATCH --gres=gpu:%d\n", config.GPU)
	content += "#SBATCH --output=result.out\n"
	content += "#SBATCH --error=error.err\n"
	content += fmt.Sprintf("#SBATCH --time=%s\n\n", config.Time)
	content += "module load gcc cuda\n\n"
	content += "make\n"
	content += fmt.Sprintf("./%s\n", config.TargetExecutableFileName)
	return content
}
