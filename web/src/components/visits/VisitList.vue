<template>
  <el-container>
    <el-main>
      <el-table
        v-loading="loading"
        :data="tableData.filter(data => !search || data.lab.toLowerCase().includes(search.toLowerCase()))"
        style="width: 100%; margin:auto"
      >
        <el-table-column type="expand">
          <template slot-scope="scope">
            <GiftListByVisit :visit-id="scope.row.id" />
            <ActivityListByVisit :visit-id="scope.row.id" />
          </template>
        </el-table-column>
        <el-table-column type="index" :index="indexMethod" />
        <el-table-column label="Lab" width="70" sortable prop="lab" />
        <el-table-column label="Arrived Date" width="200" prop="arrivedTime">
          <template slot-scope="scope">
            {{ getHumanDate(scope.row.arrivedTime) }}
          </template>
        </el-table-column>
        <el-table-column label="Departure Date" width="200" prop="departureTime">
          <template slot-scope="scope">
            {{ getHumanDate(scope.row.departureTime) }}
          </template>
        </el-table-column>
        <el-table-column label="Created By" width="120" sortable prop="createdBy" />
        <el-table-column label="Hotel Stayed" width="180" prop="hotelStayed" />
        <el-table-column label="Pickup" width="120" prop="pickup" />
        <el-table-column align="right">
          <template slot="header" slot-scope="scope">
            <el-input
              v-model="search"
              size="mini"
              placeholder="Type to search"
            />
          </template>
          <VisitUpdate
            :is-visible-update.sync="isVisibleUpdate"
            :visit.sync="visit"
            @isUpdateVisit="handleUpdate"
          />
          <VisitDelete :is-visible-delete.sync="isVisibleDelete" @isDeleteVisit="handleDelete" />
          <template slot-scope="scope">
            <el-button
              size="mini"
              @click="visit = scope.row; isVisibleUpdate = !isVisibleUpdate"
            >
              Edit
            </el-button>
            <el-button
              size="mini"
              type="danger"
              @click="isVisibleDelete = !isVisibleDelete; scopeVisit= scope"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>
</template>

<script>
import VisitUpdate from '@/components/visits/VisitUpdate.vue'
import VisitDelete from '@/components/visits/VisitDelete.vue'
// import VisitAdd from '@/components/visits/VisitAdd.vue'
import ActivityListByVisit from '@/components/activity/ActivityListByVisit.vue'
import GiftListByVisit from '@/components/gifts/GiftListByVisit.vue'
import { getHumanDate } from '@/utils/convert'
import {
  getAllVisits,
  // createVisit,
  updateVisit,
  deleteVisitById
} from '@/api/visit'

export default {
  name: 'VisitList',
  components: {
    VisitUpdate,
    VisitDelete,
    // VisitAdd,
    ActivityListByVisit,
    GiftListByVisit
  },
  data() {
    return {
      tableData: [],
      isVisibleAdd: false,
      isVisibleUpdate: false,
      isVisibleDelete: false,
      search: '',
      loading: true,
      visit: {},
      activities: {},
      scopeVisit: {}
    }
  },
  mounted() {
    getAllVisits()
      .then(resp => {
        if (resp.data != null) {
          this.tableData = resp.data
        }
        this.loading = false
      })
      .catch(err => {
        console.log(err)
      })
  },
  methods: {
    handleUpdate: function(isUpdateVisit) {
      if (isUpdateVisit) {
        this.loading = true
        updateVisit(this.visit)
          .then(resp => {
            this.$notify({
              title: 'Success',
              message: 'Update successfully!',
              type: 'success'
            })
          })
          .catch(err => {
            console.log(err)
            this.$notify.error({
              title: 'Error',
              message: err
            })
          })
        this.loading = false
      }
    },
    handleDelete: function(isDeleteVisit) {
      if (isDeleteVisit) {
        this.loading = true
        deleteVisitById(this.scopeVisit.row.id)
          .then(resp => {
            this.tableData.splice(this.scopeVisit.$index, 1)
            this.$notify({
              title: 'Success',
              message: 'Delete successfully!',
              type: 'success'
            })
          })
          .catch(err => {
            console.log(err)
            this.$notify.error({
              title: 'Error',
              message: err
            })
          })
      }
      this.loading = false
    },
    indexMethod(index) {
      return index * 1
    },
    getHumanDate
  }
}
</script>
