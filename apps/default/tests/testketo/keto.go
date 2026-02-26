package testketo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pitabwire/frame/frametests/definition"
	"github.com/pitabwire/frame/frametests/deps/testpostgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// ImageName is the Ory Keto image used for test containers.
	ImageName = "oryd/keto:latest"

	ketoConfiguration = `
version: v0.14.0

dsn: memory

serve:
  read:
    host: 0.0.0.0
    port: 4466
  write:
    host: 0.0.0.0
    port: 4467

log:
  level: debug
  format: text

namespaces:
  location: file:///home/ory/namespaces

`

	oplNamespaces = `import { Namespace, Context } from "@ory/keto-namespace-types"

class profile implements Namespace {}

class profile_tenant implements Namespace {
  related: {
    owner: profile[]
    admin: profile[]
    operator: profile[]
    viewer: profile[]

    view_profile: profile[]
    create_profile: profile[]
    update_profile: profile[]
    merge_profiles: profile[]
    manage_contacts: profile[]
    manage_roster: profile[]
    manage_relationships: profile[]
    manage_devices: profile[]
    view_devices: profile[]
    manage_geolocation: profile[]
    view_geolocation: profile[]
    ingest_location: profile[]
    manage_settings: profile[]
    view_settings: profile[]
  }

  permits = {
    view_profile: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.view_profile.includes(ctx.subject),

    create_profile: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.create_profile.includes(ctx.subject),

    update_profile: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.update_profile.includes(ctx.subject),

    merge_profiles: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.merge_profiles.includes(ctx.subject),

    manage_contacts: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_contacts.includes(ctx.subject),

    manage_roster: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_roster.includes(ctx.subject),

    manage_relationships: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_relationships.includes(ctx.subject),

    manage_devices: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_devices.includes(ctx.subject),

    view_devices: (ctx: Context): boolean =>
      this.permits.manage_devices(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.view_devices.includes(ctx.subject),

    manage_geolocation: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_geolocation.includes(ctx.subject),

    view_geolocation: (ctx: Context): boolean =>
      this.permits.manage_geolocation(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.view_geolocation.includes(ctx.subject),

    ingest_location: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.ingest_location.includes(ctx.subject),

    manage_settings: (ctx: Context): boolean =>
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.manage_settings.includes(ctx.subject),

    view_settings: (ctx: Context): boolean =>
      this.permits.manage_settings(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.view_settings.includes(ctx.subject),
  }
}
`

	namespaceFile = "/home/ory/namespaces/profile_service.ts"
)

type dependancy struct {
	*definition.DefaultImpl
}

// NewWithOpts creates a new Keto test resource with OPL namespace support.
func NewWithOpts(
	containerOpts ...definition.ContainerOption,
) definition.TestResource {
	opts := definition.ContainerOpts{
		ImageName:      ImageName,
		Ports:          []string{"4467/tcp", "4466/tcp"},
		NetworkAliases: []string{"keto", "auth-keto"},
	}
	opts.Setup(containerOpts...)

	return &dependancy{
		DefaultImpl: definition.NewDefaultImpl(opts, "http"),
	}
}

func (d *dependancy) migrateContainer(
	ctx context.Context,
	ntwk *testcontainers.DockerNetwork,
	databaseURL string,
) error {
	containerRequest := testcontainers.ContainerRequest{
		Image: d.Name(),
		Cmd:   []string{"migrate", "up", "--yes"},
		Env: map[string]string{
			"LOG_LEVEL": "debug",
			"DSN":       databaseURL,
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(ketoConfiguration),
				ContainerFilePath: "/home/ory/keto.yml",
				FileMode:          definition.ContainerFileMode,
			},
			{
				Reader:            strings.NewReader(oplNamespaces),
				ContainerFilePath: namespaceFile,
				FileMode:          definition.ContainerFileMode,
			},
		},
		WaitingFor: wait.ForExit(),
	}

	d.Configure(ctx, ntwk, &containerRequest)

	ketoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start keto migration container: %w", err)
	}

	if err = ketoContainer.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate keto migration container: %w", err)
	}
	return nil
}

func (d *dependancy) Setup(ctx context.Context, ntwk *testcontainers.DockerNetwork) error {
	if len(d.Opts().Dependencies) == 0 || !d.Opts().Dependencies[0].GetDS(ctx).IsDB() {
		return errors.New("no database dependency was supplied")
	}

	ketoDB, _, err := testpostgres.CreateDatabase(ctx, d.Opts().Dependencies[0].GetInternalDS(ctx), "keto")
	if err != nil {
		return fmt.Errorf("failed to create keto database: %w", err)
	}

	databaseURL := ketoDB.String()

	if err = d.migrateContainer(ctx, ntwk, databaseURL); err != nil {
		return err
	}

	containerRequest := testcontainers.ContainerRequest{
		Image: d.Name(),
		Cmd:   []string{"serve", "--config", "/home/ory/keto.yml"},
		Env: d.Opts().Env(map[string]string{
			"LOG_LEVEL":                 "debug",
			"LOG_LEAK_SENSITIVE_VALUES": "true",
			"DSN":                       databaseURL,
		}),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(ketoConfiguration),
				ContainerFilePath: "/home/ory/keto.yml",
				FileMode:          definition.ContainerFileMode,
			},
			{
				Reader:            strings.NewReader(oplNamespaces),
				ContainerFilePath: namespaceFile,
				FileMode:          definition.ContainerFileMode,
			},
		},
		WaitingFor: wait.ForHTTP("/health/ready").WithPort(d.DefaultPort),
	}

	d.Configure(ctx, ntwk, &containerRequest)

	ketoContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		})
	if err != nil {
		return fmt.Errorf("failed to start keto serve container: %w", err)
	}

	d.SetContainer(ketoContainer)
	return nil
}
