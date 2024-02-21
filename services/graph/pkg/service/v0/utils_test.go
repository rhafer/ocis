package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("Utils", func() {
	DescribeTable("GetDriveAndItemIDParam",
		func(driveID, itemID string, shouldPass bool) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", driveID)
			rctx.URLParams.Add("itemID", itemID)

			extractedDriveID, extractedItemID, err := service.GetDriveAndItemIDParam(
				httptest.NewRequest(http.MethodGet, "/", nil).
					WithContext(
						context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
					),
				conversions.ToPointer(log.NopLogger()),
			)

			switch shouldPass {
			case true:
				Expect(err).To(BeNil())
				parsedItemID, _ := storagespace.ParseID(itemID)
				Expect(extractedItemID).To(Equal(parsedItemID))

				parsedDriveID, _ := storagespace.ParseID(driveID)
				Expect(extractedDriveID).To(Equal(parsedDriveID))
			default:
				Expect(err).ToNot(BeNil())
			}
		},
		Entry("fails: invalid driveID", "", "1$2!3", false),
		Entry("fails: invalid itemID", "1$2", "", false),
		Entry("fails: incompatible driveID and itemID", "1$2", "3$4!5", false),
		Entry("fails: no itemID opaqueId", "1$2", "3$4", false),
		Entry("pass: valid driveID and itemID", "1$2", "1$2!5", true),
	)

	DescribeTable("IsSpaceRoot",
		func(resourceID *provider.ResourceId, isRoot bool) {
			Expect(service.IsSpaceRoot(resourceID)).To(Equal(isRoot))
		},
		Entry("spaceId and opaqueID equal", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
			SpaceId:   "2",
		}, true),
		Entry("nil", nil, false),
		Entry("spaceID empty", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
		}, false),
		Entry("opaqueID empty", &provider.ResourceId{
			StorageId: "1",
			SpaceId:   "3",
		}, false),
		Entry("spaceID and opaqueID unequal", &provider.ResourceId{
			OpaqueId: "2",
			SpaceId:  "3",
		}, false),
	)

	DescribeTable("IsShareJail",
		func(resourceID provider.ResourceId, isShareJail bool) {
			Expect(service.IsShareJail(resourceID)).To(Equal(isShareJail))
		},
		Entry("valid: share jail", provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
			SpaceId:   utils.ShareStorageSpaceID,
		}, true),
		Entry("invalid: empty storageId", provider.ResourceId{
			SpaceId: utils.ShareStorageSpaceID,
		}, false),
		Entry("invalid: empty spaceId", provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
		}, false),
		Entry("invalid: empty storageId and spaceId", provider.ResourceId{}, false),
		Entry("invalid: non share jail storageId", provider.ResourceId{
			StorageId: "123",
			SpaceId:   utils.ShareStorageSpaceID,
		}, false),
		Entry("invalid: non share jail spaceId", provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
			SpaceId:   "123",
		}, false),
		Entry("invalid: non share jail storageID and spaceId", provider.ResourceId{
			StorageId: "123",
			SpaceId:   "123",
		}, false),
	)
})
