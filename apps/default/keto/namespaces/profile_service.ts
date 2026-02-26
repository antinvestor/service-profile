import { Namespace, Context } from "@ory/keto-namespace-types"

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
