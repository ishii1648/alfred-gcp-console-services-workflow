# GCP Console Services – Alfred Workflow

A powerful workflow for quickly opening up GCP Console Services in your browser or searching for entities within them.

This workflow is inspired by [alfred-aws-console-services-workflow](https://github.com/rkoval/alfred-aws-console-services-workflow) which is a great workflow !

![GCP Console Services - Alfred Workflow Demo](alfred-gcp-console-services-workflow.gif)

## Installation

- Make sure your GCP project name are set in `~/.config/gcloud/configurations/config_default`
- [Download the latest release](https://github.com/ishii1648/alfred-gcp-console-services-workflow/releases/)
- Click to install

if you don't have `~/.config/gcloud/configurations/config_default`, please install gcloud [here](https://cloud.google.com/sdk/docs/initializing).

## Usage

To use, activate Alfred and type `gcp` to trigger this workflow. From there:

- type any search term to search for services
- if the current service result has a 🗂 in the subtitle, press Tab to autocomplete into sub-services (for example, navigate to "clusters" within the "GKE" service)
- if the current sub-service result has a 🔎 in the subtitle, press Tab again to start searching for its entities (for example, you can search for GKE clusters when tabbed to gcp GKE clusters )
