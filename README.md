# NSO Exporter

Prometheus Exporter for Tailf's NSO application

## Getting Started

Clone the repository and run the makefile

### Prerequisites

Requires atleast go 1.9, all dependencies are listed in the makefile

### Installing

The configuration file has to be placed in either /etc /home/nso_exporter or
working dir

```
git clone https://github.com/erraa/nso_exporter.git
cp example.nso_exporter.yaml nso_exporter.yaml
```

## Running the tests

make test

## Deployment

Edit the configuration file with urls, logins and which services to monitor
There is .service file included, but as of now it's manual process

## Built With

Go 1.9

## Contributing

Any contribution is appreciated.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


## Notes

### nso_exporter_commit_count

This metric is not working as intended and is not reliable.
