package models

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pitabwire/frame/security"

	geolocationv1 "buf.build/gen/go/antinvestor/geolocation/protocolbuffers/go/geolocation/v1"
)

type (
	LocationPointInput = geolocationv1.LocationPointInput
	LocationPointAPI   = geolocationv1.LocationPointObject
	AreaAPI            = geolocationv1.AreaObject
	GeoEventAPI        = geolocationv1.GeoEventObject
	NearbySubjectAPI   = geolocationv1.NearbySubjectObject
	NearbyAreaAPI      = geolocationv1.NearbyAreaObject
	AreaSubjectAPI     = geolocationv1.AreaSubjectObject
	RouteAPI           = geolocationv1.RouteObject
	RouteAssignmentAPI = geolocationv1.RouteAssignmentObject

	IngestLocationsRequest             = geolocationv1.IngestLocationsRequest
	IngestLocationsResponse            = geolocationv1.IngestLocationsResponse
	CreateAreaRequest                  = geolocationv1.CreateAreaRequest
	CreateAreaResponse                 = geolocationv1.CreateAreaResponse
	GetAreaRequest                     = geolocationv1.GetAreaRequest
	GetAreaResponse                    = geolocationv1.GetAreaResponse
	UpdateAreaRequest                  = geolocationv1.UpdateAreaRequest
	UpdateAreaResponse                 = geolocationv1.UpdateAreaResponse
	DeleteAreaRequest                  = geolocationv1.DeleteAreaRequest
	SearchAreasRequest                 = geolocationv1.SearchAreasRequest
	SearchAreasResponse                = geolocationv1.SearchAreasResponse
	CreateRouteRequest                 = geolocationv1.CreateRouteRequest
	CreateRouteResponse                = geolocationv1.CreateRouteResponse
	GetRouteRequest                    = geolocationv1.GetRouteRequest
	GetRouteResponse                   = geolocationv1.GetRouteResponse
	UpdateRouteRequest                 = geolocationv1.UpdateRouteRequest
	UpdateRouteResponse                = geolocationv1.UpdateRouteResponse
	DeleteRouteRequest                 = geolocationv1.DeleteRouteRequest
	SearchRoutesRequest                = geolocationv1.SearchRoutesRequest
	SearchRoutesResponse               = geolocationv1.SearchRoutesResponse
	AssignRouteRequest                 = geolocationv1.AssignRouteRequest
	AssignRouteResponse                = geolocationv1.AssignRouteResponse
	UnassignRouteRequest               = geolocationv1.UnassignRouteRequest
	GetSubjectRouteAssignmentsRequest  = geolocationv1.GetSubjectRouteAssignmentsRequest
	GetSubjectRouteAssignmentsResponse = geolocationv1.GetSubjectRouteAssignmentsResponse
	GetTrackRequest                    = geolocationv1.GetTrackRequest
	GetTrackResponse                   = geolocationv1.GetTrackResponse
	GetSubjectEventsRequest            = geolocationv1.GetSubjectEventsRequest
	GetSubjectEventsResponse           = geolocationv1.GetSubjectEventsResponse
	GetAreaSubjectsRequest             = geolocationv1.GetAreaSubjectsRequest
	GetAreaSubjectsResponse            = geolocationv1.GetAreaSubjectsResponse
	GetNearbySubjectsRequest           = geolocationv1.GetNearbySubjectsRequest
	GetNearbySubjectsResponse          = geolocationv1.GetNearbySubjectsResponse
	GetNearbyAreasRequest              = geolocationv1.GetNearbyAreasRequest
	GetNearbyAreasResponse             = geolocationv1.GetNearbyAreasResponse
)

type EventTenancy struct {
	TenantID    string `json:"tenant_id"`
	PartitionID string `json:"partition_id"`
	AccessID    string `json:"access_id"`
}

func ContextWithEventTenancy(
	ctx context.Context,
	tenancy EventTenancy,
	subjectID string,
) context.Context {
	if tenancy.TenantID == "" && tenancy.PartitionID == "" && tenancy.AccessID == "" {
		return ctx
	}

	claims := &security.AuthenticationClaims{
		TenantID:    tenancy.TenantID,
		PartitionID: tenancy.PartitionID,
		AccessID:    tenancy.AccessID,
		ContactID:   subjectID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subjectID,
		},
	}

	return claims.ClaimsToContext(ctx)
}

type LocationPointIngestedEvent struct {
	EventTenancy
	PointID   string  `json:"point_id"`
	SubjectID string  `json:"subject_id"`
	DeviceID  string  `json:"device_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp int64   `json:"timestamp"`
}

type AreaChangedEvent struct {
	EventTenancy
	AreaID  string `json:"area_id"`
	Action  string `json:"action"`
	OwnerID string `json:"owner_id"`
}

type RouteDeviationDetectedEvent struct {
	EventTenancy
	SubjectID      string  `json:"subject_id"`
	RouteID        string  `json:"route_id"`
	EventType      string  `json:"event_type"`
	DistanceMeters float64 `json:"distance_meters"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Timestamp      int64   `json:"timestamp"`
}

type RouteChangedEvent struct {
	EventTenancy
	RouteID string `json:"route_id"`
	Action  string `json:"action"`
	OwnerID string `json:"owner_id"`
}

type GeoEventEmitted struct {
	EventTenancy
	EventID    string       `json:"event_id"`
	SubjectID  string       `json:"subject_id"`
	AreaID     string       `json:"area_id"`
	EventType  GeoEventType `json:"event_type"`
	Timestamp  int64        `json:"timestamp"`
	Confidence float64      `json:"confidence"`
}

func ToProtoLocationSource(source LocationSource) geolocationv1.LocationSource {
	switch source {
	case LocationSourceGPS:
		return geolocationv1.LocationSource_LOCATION_SOURCE_GPS
	case LocationSourceNetwork:
		return geolocationv1.LocationSource_LOCATION_SOURCE_NETWORK
	case LocationSourceIP:
		return geolocationv1.LocationSource_LOCATION_SOURCE_IP
	case LocationSourceManual:
		return geolocationv1.LocationSource_LOCATION_SOURCE_MANUAL
	default:
		return geolocationv1.LocationSource_LOCATION_SOURCE_UNSPECIFIED
	}
}

func LocationSourceFromProto(source geolocationv1.LocationSource) LocationSource {
	switch source {
	case geolocationv1.LocationSource_LOCATION_SOURCE_NETWORK:
		return LocationSourceNetwork
	case geolocationv1.LocationSource_LOCATION_SOURCE_IP:
		return LocationSourceIP
	case geolocationv1.LocationSource_LOCATION_SOURCE_MANUAL:
		return LocationSourceManual
	case geolocationv1.LocationSource_LOCATION_SOURCE_GPS,
		geolocationv1.LocationSource_LOCATION_SOURCE_UNSPECIFIED:
		return LocationSourceGPS
	default:
		return LocationSourceGPS
	}
}

func ToProtoAreaType(areaType AreaType) geolocationv1.AreaType {
	switch areaType {
	case AreaTypeLand:
		return geolocationv1.AreaType_AREA_TYPE_LAND
	case AreaTypeBuilding:
		return geolocationv1.AreaType_AREA_TYPE_BUILDING
	case AreaTypeZone:
		return geolocationv1.AreaType_AREA_TYPE_ZONE
	case AreaTypeFence:
		return geolocationv1.AreaType_AREA_TYPE_FENCE
	case AreaTypeCustom:
		return geolocationv1.AreaType_AREA_TYPE_CUSTOM
	default:
		return geolocationv1.AreaType_AREA_TYPE_UNSPECIFIED
	}
}

func AreaTypeFromProto(areaType geolocationv1.AreaType) AreaType {
	switch areaType {
	case geolocationv1.AreaType_AREA_TYPE_BUILDING:
		return AreaTypeBuilding
	case geolocationv1.AreaType_AREA_TYPE_ZONE:
		return AreaTypeZone
	case geolocationv1.AreaType_AREA_TYPE_FENCE:
		return AreaTypeFence
	case geolocationv1.AreaType_AREA_TYPE_CUSTOM:
		return AreaTypeCustom
	case geolocationv1.AreaType_AREA_TYPE_LAND,
		geolocationv1.AreaType_AREA_TYPE_UNSPECIFIED:
		return AreaTypeLand
	default:
		return AreaTypeLand
	}
}

func ToProtoGeoEventType(eventType GeoEventType) geolocationv1.GeoEventType {
	switch eventType {
	case GeoEventTypeEnter:
		return geolocationv1.GeoEventType_GEO_EVENT_TYPE_ENTER
	case GeoEventTypeExit:
		return geolocationv1.GeoEventType_GEO_EVENT_TYPE_EXIT
	case GeoEventTypeDwell:
		return geolocationv1.GeoEventType_GEO_EVENT_TYPE_DWELL
	default:
		return geolocationv1.GeoEventType_GEO_EVENT_TYPE_UNSPECIFIED
	}
}
