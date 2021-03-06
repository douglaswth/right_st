package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/rightscale/right_st"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("RightScript Scaffold", func() {
	var (
		buffer                   *gbytes.Buffer
		tempDir                  string
		emptyScript              string
		emptyScriptMetadata      string
		metadataScript           string
		metadataScriptContents   string
		shellScript              string
		shellScriptContents      string
		shellScriptMetadata      string
		rubyScript               string
		rubyScriptContents       string
		rubyScriptMetadata       string
		perlScript               string
		perlScriptContents       string
		perlScriptMetadata       string
		powershellScript         string
		powershellScriptContents string
		powershellScriptMetadata string
	)

	BeforeEach(func() {
		buffer = gbytes.NewBuffer()
		var err error
		tempDir, err = ioutil.TempDir("", "scaffold")
		if err != nil {
			panic(err)
		}
		emptyScript = filepath.Join(tempDir, "empty.sh")
		metadataScript = filepath.Join(tempDir, "metadata.sh")
		shellScript = filepath.Join(tempDir, "shell.sh")
		rubyScript = filepath.Join(tempDir, "ruby.rb")
		perlScript = filepath.Join(tempDir, "perl.pl")
		powershellScript = filepath.Join(tempDir, "powershell.ps1")
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("With an empty script file", func() {
		BeforeEach(func() {
			emptyScriptMetadata = `# ---
# RightScript Name: Empty
# Description: (put your description here, it can be multiple lines using YAML syntax)
# Inputs: {}
# Attachments: []
# ...
`
			if err := ioutil.WriteFile(emptyScript, nil, 0600); err != nil {
				panic(err)
			}
		})

		It("should add default metadata", func() {
			err := ScaffoldRightScript(emptyScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(emptyScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(emptyScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(emptyScriptMetadata))
		})

		It("should create a backup file if desired", func() {
			err := ScaffoldRightScript(emptyScript, true, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(emptyScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(emptyScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(emptyScriptMetadata))

			info, err := os.Stat(emptyScript + ".bak")
			Expect(err).To(Succeed())
			Expect(info.Size()).To(BeEquivalentTo(0))
		})
	})

	Context("With a script with metadata", func() {
		BeforeEach(func() {
			metadataScriptContents = `#!/bin/bash
# ---
# RightScript Name: Metadata Already
# Description: A script that already has metadata
# Inputs: {}
# Attachments: []
# ...

echo 'I have metadata already!'
`
			if err := ioutil.WriteFile(metadataScript, []byte(metadataScriptContents), 0600); err != nil {
				panic(err)
			}
		})

		It("should not add metadata", func() {
			err := ScaffoldRightScript(metadataScript, false, buffer, false)
			Expect(err).To(Succeed())
			Expect(string(buffer.Contents())).Should(ContainSubstring("Script unchanged, already contains metadata"))

			script, err := ioutil.ReadFile(metadataScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(metadataScriptContents))
		})
	})

	Context("With rescaffolding a script with metadata", func() {
		var metadataScriptBefore string
		var metadataScriptAfter string
		BeforeEach(func() {
			metadataScriptBefore = `#!/bin/bash
# ---
# RightScript Name: Metadata Already
# Description: A script that already has metadata
# Inputs:
#   FOO:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: true
#     Advanced: true
#   BAR:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
# Attachments: []
# ...

echo $FOO
# BAR was removed
${BAZ:=hello}
echo $BAZ
`
			metadataScriptAfter = `#!/bin/bash
# ---
# RightScript Name: Metadata Already
# Description: A script that already has metadata
# Inputs:
#   FOO:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: true
#     Advanced: true
#   BAZ:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
#     Default: text:hello
# Attachments: []
# ...

echo $FOO
# BAR was removed
${BAZ:=hello}
echo $BAZ
`
			if err := ioutil.WriteFile(metadataScript, []byte(metadataScriptBefore), 0600); err != nil {
				panic(err)
			}
		})

		It("should re-scaffold metadata", func() {
			err := ScaffoldRightScript(metadataScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(metadataScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(metadataScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(metadataScriptAfter))
		})
	})

	Context("With a shell script", func() {
		BeforeEach(func() {
			shebang := "#!/bin/bash\n"
			shellScriptContents = `
: ${STRING:=hello}
: ${ARRAY:=hello,world}

echo "$STRING $ARRAY $PATH"
echo $RS_ATTACH_DIR/attachment.txt
`
			shellScriptMetadata = shebang + `# ---
# RightScript Name: Shell
# Description: (put your description here, it can be multiple lines using YAML syntax)
# Inputs:
#   STRING:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
#     Default: text:hello
#   ARRAY:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: array
#     Required: false
#     Advanced: false
#     Default: array:["text:hello","text:world"]
# Attachments:
# - attachment.txt
# ...
` + shellScriptContents
			shellScriptContents = shebang + shellScriptContents
			if err := ioutil.WriteFile(shellScript, []byte(shellScriptContents), 0600); err != nil {
				panic(err)
			}
		})

		It("should add metadata with variables and their default values", func() {
			err := ScaffoldRightScript(shellScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(shellScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(shellScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(shellScriptMetadata))
		})
	})

	Context("With a Ruby script", func() {
		BeforeEach(func() {
			shebang := "#!/usr/bin/env ruby\n"
			rubyScriptContents = `
puts "#{ENV['INPUT']} #{ENV["PATH"]}"
puts "#{ENV['RS_ATTACH_DIR']}/attachment.txt"
`
			rubyScriptMetadata = shebang + `# ---
# RightScript Name: Ruby
# Description: (put your description here, it can be multiple lines using YAML syntax)
# Inputs:
#   INPUT:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
# Attachments:
# - attachment.txt
# ...
` + rubyScriptContents
			rubyScriptContents = shebang + rubyScriptContents
			if err := ioutil.WriteFile(rubyScript, []byte(rubyScriptContents), 0600); err != nil {
				panic(err)
			}
		})

		It("should add metadata with variables", func() {
			err := ScaffoldRightScript(rubyScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(rubyScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(rubyScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(rubyScriptMetadata))
		})
	})

	Context("With a Perl script", func() {
		BeforeEach(func() {
			shebang := "#!/usr/bin/env perl\n"
			perlScriptContents = `
print "$ENV{INPUT} $ENV{PATH}\n";
print "$ENV{RS_ATTACH_DIR}/attachment.txt" . "\n";
`
			perlScriptMetadata = shebang + `# ---
# RightScript Name: Perl
# Description: (put your description here, it can be multiple lines using YAML syntax)
# Inputs:
#   INPUT:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
# Attachments:
# - attachment.txt
# ...
` + perlScriptContents
			perlScriptContents = shebang + perlScriptContents
			if err := ioutil.WriteFile(perlScript, []byte(perlScriptContents), 0600); err != nil {
				panic(err)
			}
		})

		It("should add metadata with variables", func() {
			err := ScaffoldRightScript(perlScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(perlScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(perlScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(perlScriptMetadata))
		})
	})

	Context("With a PowerShell script", func() {
		BeforeEach(func() {
			powershellScriptContents = `
Write-Output "${env:INPUT} $env:PATH"
Write-Output "${env:RS_ATTACH_DIR}/attachment.txt"
`
			powershellScriptMetadata = `# ---
# RightScript Name: Powershell
# Description: (put your description here, it can be multiple lines using YAML syntax)
# Inputs:
#   INPUT:
#     Category: (put your input category here)
#     Description: (put your input description here, it can be multiple lines using
#       YAML syntax)
#     Input Type: single
#     Required: false
#     Advanced: false
# Attachments:
# - attachment.txt
# ...
` + powershellScriptContents
			if err := ioutil.WriteFile(powershellScript, []byte(powershellScriptContents), 0600); err != nil {
				panic(err)
			}
		})

		It("should add metadata with variables", func() {
			err := ScaffoldRightScript(powershellScript, false, buffer, true)
			Expect(err).To(Succeed())
			Expect(buffer.Contents()).To(BeEquivalentTo(powershellScript + ": Added metadata\n"))

			script, err := ioutil.ReadFile(powershellScript)
			Expect(err).To(Succeed())
			Expect(script).To(BeEquivalentTo(powershellScriptMetadata))
		})
	})
})
