package redis_test

import (
	"context"
	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/redis/go-redis/v9"
)

var _ = Describe("Probabilistic commands", Label("probabilistic"), func() {
	ctx := context.TODO()
	var client *redis.Client

	BeforeEach(func() {
		client = redis.NewClient(redisOptions())
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	Describe("bloom", Label("bloom"), func() {
		It("should BFAdd", Label("bloom", "bfadd"), func() {
			resultAdd, err := client.BFAdd(ctx, "testbf1", 1).Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(resultAdd).To(BeTrue())

			resultInfo, err := client.BFInfo(ctx, "testbf1").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(resultInfo).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(resultInfo.NumItemsInserted).To(BeEquivalentTo(int64(1)))
		})

		It("should BFCard", Label("bloom", "bfcard"), func() {
			// This is a probabilistic data structure, and it's not always guaranteed that we will get back
			// the exact number of inserted items, during hash collisions
			// But with such a low number of items (only 3),
			// the probability of a collision is very low, so we can expect to get back the exact number of items
			_, err := client.BFAdd(ctx, "testbf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			_, err = client.BFAdd(ctx, "testbf1", "item2").Result()
			Expect(err).NotTo(HaveOccurred())
			_, err = client.BFAdd(ctx, "testbf1", 3).Result()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.BFCard(ctx, "testbf1").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEquivalentTo(int64(3)))
		})

		It("should BFExists", Label("bloom", "bfexists"), func() {
			exists, err := client.BFExists(ctx, "testbf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())

			_, err = client.BFAdd(ctx, "testbf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())

			exists, err = client.BFExists(ctx, "testbf1", "item1").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should BFInfo and BFReserve", Label("bloom", "bfinfo", "bfreserve"), func() {
			err := client.BFReserve(ctx, "testbf1", 0.001, 2000).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.BFInfo(ctx, "testbf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(result.Capacity).To(BeEquivalentTo(int64(2000)))
		})

		It("should BFInfoArg", Label("bloom", "bfinfoarg"), func() {
			err := client.BFReserve(ctx, "testbf1", 0.001, 2000).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.BFInfoArg(ctx, "testbf1", redis.BFInfoArgs(redis.BFCAPACITY)).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Capacity).To(BeEquivalentTo(int64(2000)))
		})

		It("should BFInsert", Label("bloom", "bfinsert"), func() {
			options := &redis.BFReserveOptions{
				Capacity:   2000,
				Error:      0.001,
				Expansion:  3,
				NonScaling: false,
			}
			resultInsert, err := client.BFInsert(ctx, "testbf1", options, "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultInsert)).To(BeEquivalentTo(1))

			exists, err := client.BFExists(ctx, "testbf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			result, err := client.BFInfo(ctx, "testbf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(result.Capacity).To(BeEquivalentTo(int64(2000)))
			Expect(result.ExpansionRate).To(BeEquivalentTo(int64(3)))
		})

		It("should BFMAdd", Label("bloom", "bfmadd"), func() {
			resultAdd, err := client.BFMAdd(ctx, "testbf1", "item1", "item2", "item3").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultAdd)).To(Equal(3))

			resultInfo, err := client.BFInfo(ctx, "testbf1").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(resultInfo).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(resultInfo.NumItemsInserted).To(BeEquivalentTo(int64(3)))
		})

		It("should BFMExists", Label("bloom", "bfmexists"), func() {
			exist, err := client.BFMExists(ctx, "testbf1", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(exist)).To(Equal(int(3)))
			Expect(exist[0]).To(BeFalse())
			Expect(exist[1]).To(BeFalse())
			Expect(exist[2]).To(BeFalse())

			_, err = client.BFMAdd(ctx, "testbf1", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())

			exist, err = client.BFMExists(ctx, "testbf1", "item1", "item2", "item3", "item4").Result()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(exist)).To(Equal(int(4)))
			Expect(exist[0]).To(BeTrue())
			Expect(exist[1]).To(BeTrue())
			Expect(exist[2]).To(BeTrue())
			Expect(exist[3]).To(BeFalse())
		})

		It("should BFReserveExpansion", Label("bloom", "bfreserveexpansion"), func() {
			err := client.BFReserveExpansion(ctx, "testbf1", 0.001, 2000, 3).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.BFInfo(ctx, "testbf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(result.Capacity).To(BeEquivalentTo(int64(2000)))
			Expect(result.ExpansionRate).To(BeEquivalentTo(int64(3)))
		})

		It("should BFReserveArgs", Label("bloom", "bfreserveargs"), func() {
			options := &redis.BFReserveOptions{
				Capacity:   2000,
				Error:      0.001,
				Expansion:  3,
				NonScaling: false,
			}
			err := client.BFReserveArgs(ctx, "testbf", options).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.BFInfo(ctx, "testbf").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.BFInfo{}))
			Expect(result.Capacity).To(BeEquivalentTo(int64(2000)))
			Expect(result.ExpansionRate).To(BeEquivalentTo(int64(3)))
		})
	})

	Describe("cuckoo", Label("cuckoo"), func() {
		It("should CFAdd", Label("cuckoo", "cfadd"), func() {
			add, err := client.CFAdd(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(add).To(BeTrue())

			exists, err := client.CFExists(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			info, err := client.CFInfo(ctx, "testcf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(BeAssignableToTypeOf(redis.CFInfo{}))
			Expect(info.NumItemsInserted).To(BeEquivalentTo(int64(1)))
		})

		It("should CFAddNX", Label("cuckoo", "cfaddnx"), func() {
			add, err := client.CFAddNX(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(add).To(BeTrue())

			exists, err := client.CFExists(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			result, err := client.CFAddNX(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeFalse())

			info, err := client.CFInfo(ctx, "testcf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(BeAssignableToTypeOf(redis.CFInfo{}))
			Expect(info.NumItemsInserted).To(BeEquivalentTo(int64(1)))
		})

		It("should CFCount", Label("cuckoo", "cfcount"), func() {
			err := client.CFAdd(ctx, "testcf1", "item1").Err()
			cnt, err := client.CFCount(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(BeEquivalentTo(int64(1)))

			err = client.CFAdd(ctx, "testcf1", "item1").Err()
			Expect(err).NotTo(HaveOccurred())

			cnt, err = client.CFCount(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(BeEquivalentTo(int64(2)))
		})

		It("should CFDel and CFExists", Label("cuckoo", "cfdel", "cfexists"), func() {
			err := client.CFAdd(ctx, "testcf1", "item1").Err()
			Expect(err).NotTo(HaveOccurred())

			exists, err := client.CFExists(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			del, err := client.CFDel(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(del).To(BeTrue())

			exists, err = client.CFExists(ctx, "testcf1", "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should CFInfo and CFReserve", Label("cuckoo", "cfinfo", "cfreserve"), func() {
			err := client.CFReserve(ctx, "testcf1", 1000).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.CFInfo(ctx, "testcf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.CFInfo{}))
		})

		It("should CFInfo and CFReserveArgs", Label("cuckoo", "cfinfo", "cfreserveargs"), func() {
			args := &redis.CFReserveOptions{
				Capacity:      2048,
				BucketSize:    3,
				MaxIterations: 15,
				Expansion:     2,
			}

			err := client.CFReserveArgs(ctx, "testcf1", args).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.CFInfo(ctx, "testcf1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeAssignableToTypeOf(redis.CFInfo{}))
			Expect(result.BucketSize).To(BeEquivalentTo(int64(3)))
			Expect(result.MaxIteration).To(BeEquivalentTo(int64(15)))
			Expect(result.ExpansionRate).To(BeEquivalentTo(int64(2)))
		})

		It("should CFInsert", Label("cuckoo", "cfinsert"), func() {
			args := &redis.CFInsertOptions{
				Capacity: 3000,
				NoCreate: true,
			}

			result, err := client.CFInsert(ctx, "testcf1", args, "item1", "item2", "item3").Result()
			Expect(err).To(HaveOccurred())

			args = &redis.CFInsertOptions{
				Capacity: 3000,
				NoCreate: false,
			}

			result, err = client.CFInsert(ctx, "testcf1", args, "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(3))
		})

		It("should CFInsertNx", Label("cuckoo", "cfinsertnx"), func() {
			args := &redis.CFInsertOptions{
				Capacity: 3000,
				NoCreate: true,
			}

			result, err := client.CFInsertNx(ctx, "testcf1", args, "item1", "item2", "item2").Result()
			Expect(err).To(HaveOccurred())

			args = &redis.CFInsertOptions{
				Capacity: 3000,
				NoCreate: false,
			}

			result, err = client.CFInsertNx(ctx, "testcf2", args, "item1", "item2", "item2").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(3))
			Expect(result[0]).To(BeEquivalentTo(int64(1)))
			Expect(result[1]).To(BeEquivalentTo(int64(1)))
			Expect(result[2]).To(BeEquivalentTo(int64(0)))
		})

		It("should CFMexists", Label("cuckoo", "cfmexists"), func() {
			err := client.CFInsert(ctx, "testcf1", nil, "item1", "item2", "item3").Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.CFMExists(ctx, "testcf1", "item1", "item2", "item3", "item4").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(4))
			Expect(result[0]).To(BeTrue())
			Expect(result[1]).To(BeTrue())
			Expect(result[2]).To(BeTrue())
			Expect(result[3]).To(BeFalse())
		})

	})

	Describe("CMS", Label("cms"), func() {
		It("should CMSIncrBy", Label("cms", "cmsincrby"), func() {
			err := client.CMSInitByDim(ctx, "testcms1", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.CMSIncrBy(ctx, "testcms1", "item1", 1, "item2", 2, "item3", 3).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(3))
			Expect(result[0]).To(BeEquivalentTo(int64(1)))
			Expect(result[1]).To(BeEquivalentTo(int64(2)))
			Expect(result[2]).To(BeEquivalentTo(int64(3)))

		})

		It("should CMSInitByDim and CMSInfo", Label("cms", "cmsinitbydim", "cmsinfo"), func() {
			err := client.CMSInitByDim(ctx, "testcms1", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())

			info, err := client.CMSInfo(ctx, "testcms1").Result()
			Expect(err).NotTo(HaveOccurred())

			Expect(info).To(BeAssignableToTypeOf(redis.CMSInfo{}))
			Expect(info.Width).To(BeEquivalentTo(int64(5)))
			Expect(info.Depth).To(BeEquivalentTo(int64(10)))
		})

		It("should CMSInitByProb", Label("cms", "cmsinitbyprob"), func() {
			err := client.CMSInitByProb(ctx, "testcms1", 0.002, 0.01).Err()
			Expect(err).NotTo(HaveOccurred())

			info, err := client.CMSInfo(ctx, "testcms1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(BeAssignableToTypeOf(redis.CMSInfo{}))
		})

		It("should CMSMerge, CMSMergeWithWeight and CMSQuery", Label("cms", "cmsmerge", "cmsquery"), func() {
			err := client.CMSMerge(ctx, "destCms1", "testcms2", "testcms3").Err()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("CMS: key does not exist"))

			err = client.CMSInitByDim(ctx, "destCms1", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.CMSInitByDim(ctx, "destCms2", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.CMSInitByDim(ctx, "cms1", 2, 20).Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.CMSInitByDim(ctx, "cms2", 3, 20).Err()
			Expect(err).NotTo(HaveOccurred())

			err = client.CMSMerge(ctx, "destCms1", "cms1", "cms2").Err()
			Expect(err).To(MatchError("CMS: width/depth is not equal"))

			client.Del(ctx, "cms1", "cms2")

			err = client.CMSInitByDim(ctx, "cms1", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())
			err = client.CMSInitByDim(ctx, "cms2", 5, 10).Err()
			Expect(err).NotTo(HaveOccurred())

			client.CMSIncrBy(ctx, "cms1", "item1", 1, "item2", 2)
			client.CMSIncrBy(ctx, "cms2", "item2", 2, "item3", 3)

			err = client.CMSMerge(ctx, "destCms1", "cms1", "cms2").Err()
			Expect(err).NotTo(HaveOccurred())

			result, err := client.CMSQuery(ctx, "destCms1", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(3))
			Expect(result[0]).To(BeEquivalentTo(int64(1)))
			Expect(result[1]).To(BeEquivalentTo(int64(4)))
			Expect(result[2]).To(BeEquivalentTo(int64(3)))

			sourceSketches := map[string]int{
				"cms1": 1,
				"cms2": 2,
			}
			err = client.CMSMergeWithWeight(ctx, "destCms2", sourceSketches).Err()
			Expect(err).NotTo(HaveOccurred())

			result, err = client.CMSQuery(ctx, "destCms2", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(3))
			Expect(result[0]).To(BeEquivalentTo(int64(1)))
			Expect(result[1]).To(BeEquivalentTo(int64(6)))
			Expect(result[2]).To(BeEquivalentTo(int64(6)))

		})

	})

	Describe("TopK", Label("topk"), func() {
		It("should TopKReserve, TOPKInfo, TOPKAdd, TOPKQuery, TOPKCount, TOPKIncrBy, TOPKList, TOPKListWithCount", Label("topk", "topkreserve", "topkinfo", "topkadd", "topkquery", "topkcount", "topkincrby", "topklist", "topklistwithcount"), func() {
			err := client.TOPKReserve(ctx, "topk1", 3).Err()
			Expect(err).NotTo(HaveOccurred())

			resultInfo, err := client.TOPKInfo(ctx, "topk1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(resultInfo.K).To(BeEquivalentTo(int64(3)))

			resultAdd, err := client.TOPKAdd(ctx, "topk1", "item1", "item2", 3, "item1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultAdd)).To(BeEquivalentTo(int64(4)))

			resultQuery, err := client.TOPKQuery(ctx, "topk1", "item1", "item2", 4, 3).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultQuery)).To(BeEquivalentTo(int(4)))
			Expect(resultQuery[0]).To(BeTrue())
			Expect(resultQuery[1]).To(BeTrue())
			Expect(resultQuery[2]).To(BeFalse())
			Expect(resultQuery[3]).To(BeTrue())

			resultCount, err := client.TOPKCount(ctx, "topk1", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultCount)).To(BeEquivalentTo(int(3)))
			Expect(resultCount[0]).To(BeEquivalentTo(int64(2)))
			Expect(resultCount[1]).To(BeEquivalentTo(int64(1)))
			Expect(resultCount[2]).To(BeEquivalentTo(int64(0)))

			resultIncr, err := client.TOPKIncrBy(ctx, "topk1", "item1", 5, "item2", 10).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultIncr)).To(BeEquivalentTo(int(2)))

			resultCount, err = client.TOPKCount(ctx, "topk1", "item1", "item2", "item3").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultCount)).To(BeEquivalentTo(int(3)))
			Expect(resultCount[0]).To(BeEquivalentTo(int64(7)))
			Expect(resultCount[1]).To(BeEquivalentTo(int64(11)))
			Expect(resultCount[2]).To(BeEquivalentTo(int64(0)))

			resultList, err := client.TOPKList(ctx, "topk1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultList)).To(BeEquivalentTo(int(3)))
			Expect(resultList).To(ContainElements("item2", "item1", "3"))

			resultListWithCount, err := client.TOPKListWithCount(ctx, "topk1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultListWithCount)).To(BeEquivalentTo(int(3)))
			Expect(resultListWithCount["3"]).To(BeEquivalentTo(int64(1)))
			Expect(resultListWithCount["item1"]).To(BeEquivalentTo(int64(7)))
			Expect(resultListWithCount["item2"]).To(BeEquivalentTo(int64(11)))
		})

		It("should TOPKReserveWithOptions", Label("topk", "topkreservewithoptions"), func() {
			//&redis.TOPKReserveOptions{Width: 1500, Depth: 8, Decay: 0.5}
			err := client.TOPKReserveWithOptions(ctx, "topk1", 3, 1500, 8, 0.5).Err()
			Expect(err).NotTo(HaveOccurred())

			resultInfo, err := client.TOPKInfo(ctx, "topk1").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(resultInfo.K).To(BeEquivalentTo(int64(3)))
			Expect(resultInfo.Width).To(BeEquivalentTo(int64(1500)))
			Expect(resultInfo.Depth).To(BeEquivalentTo(int64(8)))
			Expect(resultInfo.Decay).To(BeEquivalentTo(float64(0.5)))
		})

	})
})