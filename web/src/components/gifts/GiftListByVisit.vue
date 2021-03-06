<template>
  <el-container style="background: #f0f9eb;">
    <el-header style="width: 100%; margin:auto">
      <el-button
        type="success"
        icon="el-icon-plus"
        plain
        style="float:right"
        @click="getGiftAssociate"
      >
        Assign gift
      </el-button>
    </el-header>
    <el-main>
      <el-table
        v-loading="loading"
        :data="tableData.filter(data => !search || data.name.toLowerCase().includes(search.toLowerCase()))"
        style="width: 100%"
      >
        <el-table-column label="Name" prop="giftName" sortable />
        <el-table-column label="Quantity" prop="quantity" sortable />
        <el-table-column label="Note" prop="note" sortable />
        <el-table-column align="right">
          <GiftAssociateUpdate
            :is-visible-update.sync="isVisibleUpdate"
            :gift.sync="gift"
            @isUpdateGift="handleGiftAssociateUpdate"
          />
          <GiftAssociateDelete
            :is-visible-delete.sync="isVisibleDelete"
            :gift-name.sync="giftName"
            @isDeleteGift="handleGiftAssociateDelete"
          />
          <GiftAssociateAdd :assigned-gifts.sync="assignedGifts" :is-visible-assign.sync="isVisibleAssign" @isGiftAssociateAdd="handleGiftAssociateAdd" />
          <template slot-scope="scope">
            <el-button
              size="mini"
              @click="gift = scope.row; isVisibleUpdate = !isVisibleUpdate; "
            >
              Edit
            </el-button>
            <el-button
              size="mini"
              type="danger"
              @click="isVisibleDelete = !isVisibleDelete; scopeGift = scope; giftName = scope.row.giftName"
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
import GiftAssociateUpdate from '@/components/giftAssocicates/GiftAssociateUpdate.vue'
import GiftAssociateDelete from '@/components/giftAssocicates/GiftAssociateDelete.vue'
import GiftAssociateAdd from '@/components/giftAssocicates/GiftAssociateAdd.vue'

import {
  getGiftAssociatesByVisitID,
  createGiftAssociate,
  deleteGiftAssociateById,
  modifyGiftAssociates
} from '@/api/giftAssociate'

export default {
  name: 'GiftListByVisit',
  components: {
    GiftAssociateUpdate,
    GiftAssociateDelete,
    GiftAssociateAdd
  },
  props: {
    visitId: { type: String, default: '' },
    friendId: { type: String, default: '' },
    friendName: { type: String, default: '' }
  },

  data() {
    return {
      tableData: [],
      isVisibleUpdate: false,
      isVisibleDelete: false,
      isVisibleAssign: false,
      search: '',
      giftAssociates: {},
      giftName: '',
      loading: true,
      scopeGift: {},
      gift: {},
      assignedGifts: []
    }
  },
  mounted() {
    getGiftAssociatesByVisitID(this.visitId)
      .then(resp => {
        if (resp.data != null) {
          this.tableData = resp.data
          this.tableData.forEach((gift, index) => {
            this.assignedGifts.push(gift.id)
          })
        }
        this.loading = false
      })
      .catch(err => {
        console.log(err)
      })
  },

  methods: {
    getGiftAssociate: function() {
      this.assignedGifts = []
      this.tableData.forEach((gift, index) => {
        this.assignedGifts.push(gift.giftId)
      })
      this.isVisibleAssign = !this.isVisibleAssign
    },

    // handleGiftAssociateAdd: Create if this gift is not exist in current assigned gifts.
    // Delete if the current assigned gifts does not exist in updated assigned gifts.
    handleGiftAssociateAdd: function(isGiftAssociateAdd, updatedGiftAssociates) {
      if (isGiftAssociateAdd) {
        var currentAssignedGifts = this.assignedGifts
        updatedGiftAssociates.forEach((gift, index) => {
          // If currentAssignedGifts doesn not include this gift, assign this gift
          var position = currentAssignedGifts.indexOf(gift.initial)
          if (position < 0) {
            var giftAssociate = {}
            giftAssociate.giftId = gift.initial
            giftAssociate.visitId = this.visitId
            giftAssociate.friendId = this.friendId
            giftAssociate.giftName = gift.label
            giftAssociate.friendName = this.friendName
            giftAssociate.quantity = 1
            giftAssociate.note = ''
            createGiftAssociate(giftAssociate)
              .then(resp => {
                this.$notify({
                  title: 'Success',
                  message: 'Update successfully!',
                  type: 'success'
                })
                giftAssociate.id = resp.data.id
                this.tableData.splice(0, 0, giftAssociate)
              })
              .catch(err => {
                console.log(err)
                this.$notify.error({
                  title: 'Error',
                  message: err
                })
              })
          } else {
            // Remove this gifts are not in currentAssignedGifts
            currentAssignedGifts.splice(position, 1)
          }
        })
        // Current elements in currentAssignedGifts are gifts that is not in updated Gift
        if (currentAssignedGifts.length > 0) {
          currentAssignedGifts.forEach((id) => {
            // Find to get gift associate ID, then delete this gift associate
            var data = this.tableData.filter(data => !id || data.giftId.includes(id))
            deleteGiftAssociateById(data[0].id)
              .then(resp => {
                var giftIndex = this.assignedGifts.indexOf(id)
                this.tableData.splice(giftIndex, 1)
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
          })
        }
        this.loading = false
      }
    },

    handleGiftAssociateUpdate: function(isUpdateGift) {
      if (isUpdateGift) {
        this.loading = true
        modifyGiftAssociates(this.gift)
          .then(resp => {
            console.log(resp.data)
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
    handleGiftAssociateDelete: function(isDeleteGift) {
      if (isDeleteGift) {
        this.loading = true
        deleteGiftAssociateById(this.scopeGift.row.id)
          .then(resp => {
            this.tableData.splice(this.scopeGift.$index, 1)
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
    }
  }
}
</script>
