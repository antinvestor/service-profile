import { Namespace, Context } from "@ory/keto-namespace-types"

class profile_user implements Namespace {}

class tenancy_access implements Namespace {
  related: {
    member: (profile_user | tenancy_access)[]
    service: profile_user[]
  }
}

class service_profile implements Namespace {
  related: {
    owner: profile_user[]
    admin: profile_user[]
    operator: profile_user[]
    viewer: profile_user[]
    member: profile_user[]
    service: (profile_user | tenancy_access)[]

    granted_profile_view: (profile_user | service_profile)[]
    granted_profile_create: (profile_user | service_profile)[]
    granted_profile_update: (profile_user | service_profile)[]
    granted_profiles_merge: (profile_user | service_profile)[]
    granted_contacts_manage: (profile_user | service_profile)[]
    granted_roster_manage: (profile_user | service_profile)[]
    granted_relationships_manage: (profile_user | service_profile)[]
    granted_devices_manage: (profile_user | service_profile)[]
    granted_devices_view: (profile_user | service_profile)[]
    granted_geolocation_manage: (profile_user | service_profile)[]
    granted_geolocation_view: (profile_user | service_profile)[]
    granted_location_ingest: (profile_user | service_profile)[]
    granted_settings_manage: (profile_user | service_profile)[]
    granted_settings_view: (profile_user | service_profile)[]
  }

  permits = {
    profile_view: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.member.includes(ctx.subject) ||
      this.related.granted_profile_view.includes(ctx.subject),

    profile_create: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_profile_create.includes(ctx.subject),

    profile_update: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_profile_update.includes(ctx.subject),

    profiles_merge: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_profiles_merge.includes(ctx.subject),

    contacts_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_contacts_manage.includes(ctx.subject),

    roster_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_roster_manage.includes(ctx.subject),

    relationships_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_relationships_manage.includes(ctx.subject),

    devices_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_devices_manage.includes(ctx.subject),

    devices_view: (ctx: Context): boolean =>
      this.permits.devices_manage(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.member.includes(ctx.subject) ||
      this.related.granted_devices_view.includes(ctx.subject),

    geolocation_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_geolocation_manage.includes(ctx.subject),

    geolocation_view: (ctx: Context): boolean =>
      this.permits.geolocation_manage(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.member.includes(ctx.subject) ||
      this.related.granted_geolocation_view.includes(ctx.subject),

    location_ingest: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.granted_location_ingest.includes(ctx.subject),

    settings_manage: (ctx: Context): boolean =>
      this.related.service.includes(ctx.subject) ||
      this.related.owner.includes(ctx.subject) ||
      this.related.admin.includes(ctx.subject) ||
      this.related.granted_settings_manage.includes(ctx.subject),

    settings_view: (ctx: Context): boolean =>
      this.permits.settings_manage(ctx) ||
      this.related.operator.includes(ctx.subject) ||
      this.related.viewer.includes(ctx.subject) ||
      this.related.member.includes(ctx.subject) ||
      this.related.granted_settings_view.includes(ctx.subject),
  }
}
