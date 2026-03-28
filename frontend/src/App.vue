<template>
  <div class="container mt-4">
    <h1 class="mb-4">Wallet Manager</h1>
    
    <div class="card mb-4">
      <div class="card-header">
        <h5 class="mb-0">Создать новый кошелек</h5>
      </div>
      <div class="card-body">
        <button @click="createWallet" class="btn btn-success" :disabled="loading">
          <i class="fas fa-plus"></i> Создать кошелек
        </button>
      </div>
    </div>
    
    <div class="card mb-4">
      <div class="card-header">
        <h5 class="mb-0">Список кошельков (первые 10)</h5>
      </div>
      <div class="card-body">
        <button @click="loadWallets" class="btn btn-primary mb-3" :disabled="loading">
          <i class="fas fa-sync-alt"></i> Обновить список
        </button>
        
        <div v-if="loading" class="text-center">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Загрузка...</span>
          </div>
        </div>
        
        <div v-else-if="wallets.length === 0" class="alert alert-info">
          Нет кошельков. Создайте первый!
        </div>
        
        <div v-else class="table-responsive">
          <table class="table table-hover">
            <thead>
              <tr>
                <th>ID</th>
                <th>Баланс</th>
                <th>Создан</th>
                <th>Обновлен</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wallet in wallets" :key="wallet.id">
                <td><code>{{ wallet.id.substring(0, 8) }}...</code></td>
                <td><strong>{{ wallet.amount }} ₽</strong></td>
                <td>{{ formatDate(wallet.createdAt) }}</td>
                <td>{{ formatDate(wallet.updatedAt) }}</td>
                <td>
                  <button @click="selectWallet(wallet)" class="btn btn-sm btn-info me-1" data-bs-toggle="modal" data-bs-target="#walletModal">
                    <i class="fas fa-exchange-alt"></i>
                  </button>
                  <button @click="deleteWallet(wallet.id)" class="btn btn-sm btn-danger" :disabled="loading">
                    <i class="fas fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- Модальное окно -->
    <div class="modal fade" id="walletModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content" v-if="selectedWallet">
          <div class="modal-header">
            <h5 class="modal-title">Кошелек: {{ selectedWallet.id }}</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label"><strong>Текущий баланс:</strong></label>
              <h3>{{ selectedWallet.amount }} ₽</h3>
            </div>
            
            <div class="mb-3">
              <label class="form-label"><strong>Тип операции:</strong></label>
              <select v-model="operationType" class="form-select">
                <option value="deposit">Пополнение (deposit)</option>
                <option value="withdraw">Снятие (withdraw)</option>
              </select>
            </div>
            
            <div class="mb-3">
              <label class="form-label"><strong>Сумма:</strong></label>
              <input type="number" v-model="operationAmount" class="form-control" placeholder="Введите сумму">
            </div>
            
            <div v-if="errorMessage" class="alert alert-danger">{{ errorMessage }}</div>
            <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
            <button @click="changeBalance" class="btn btn-primary" :disabled="loading">
              <i class="fas fa-exchange-alt"></i> Выполнить операцию
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { Modal } from 'bootstrap'

// API base URL - используем прокси
const API_BASE_URL = '/api/v1'

const wallets = ref([])
const selectedWallet = ref(null)
const operationType = ref('deposit')
const operationAmount = ref(null)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const loadWallets = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await axios.get(`${API_BASE_URL}/wallet`)
    wallets.value = response.data
  } catch (error) {
    errorMessage.value = `Ошибка загрузки: ${error.response?.data?.message || error.message}`
    console.error(error)
  } finally {
    loading.value = false
  }
}

const createWallet = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await axios.post(`${API_BASE_URL}/wallet`)
    wallets.value.unshift(response.data)
    alert(`Кошелек создан! ID: ${response.data.id}`)
  } catch (error) {
    errorMessage.value = `Ошибка создания: ${error.response?.data?.message || error.message}`
    alert(errorMessage.value)
  } finally {
    loading.value = false
  }
}

const selectWallet = (wallet) => {
  selectedWallet.value = { ...wallet }
  operationType.value = 'deposit'
  operationAmount.value = null
  errorMessage.value = ''
  successMessage.value = ''
}

const changeBalance = async () => {
  if (!operationAmount.value || operationAmount.value <= 0) {
    errorMessage.value = 'Введите корректную сумму (больше 0)'
    return
  }
  
  if (operationType.value === 'withdraw' && operationAmount.value > selectedWallet.value.amount) {
    errorMessage.value = 'Недостаточно средств на счете'
    return
  }
  
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const payload = {
      walletId: selectedWallet.value.id,
      operationType: operationType.value,
      amount: parseInt(operationAmount.value)
    }
    
    await axios.post(`${API_BASE_URL}/wallet/change-balance`, payload)
    
    const newAmount = operationType.value === 'deposit' 
      ? selectedWallet.value.amount + parseInt(operationAmount.value)
      : selectedWallet.value.amount - parseInt(operationAmount.value)
    
    selectedWallet.value.amount = newAmount
    
    const index = wallets.value.findIndex(w => w.id === selectedWallet.value.id)
    if (index !== -1) {
      wallets.value[index].amount = newAmount
      wallets.value[index].updatedAt = new Date().toISOString()
    }
    
    successMessage.value = `Операция выполнена! Новый баланс: ${newAmount} ₽`
    operationAmount.value = null
    
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (error) {
    errorMessage.value = `Ошибка операции: ${error.response?.data?.message || error.message}`
    console.error(error)
  } finally {
    loading.value = false
  }
}

const deleteWallet = async (walletId) => {
  if (!confirm('Вы уверены, что хотите удалить этот кошелек?')) {
    return
  }
  
  loading.value = true
  errorMessage.value = ''
  
  try {
    await axios.delete(`${API_BASE_URL}/wallet/${walletId}`)
    wallets.value = wallets.value.filter(w => w.id !== walletId)
    
    if (selectedWallet.value && selectedWallet.value.id === walletId) {
      selectedWallet.value = null
      const modal = Modal.getInstance(document.getElementById('walletModal'))
      if (modal) modal.hide()
    }
    
    alert('Кошелек удален')
  } catch (error) {
    errorMessage.value = `Ошибка удаления: ${error.response?.data?.message || error.message}`
    alert(errorMessage.value)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '—'
  const date = new Date(dateString)
  return date.toLocaleString('ru-RU')
}

onMounted(() => {
  loadWallets()
})
</script>

<style scoped>
/* Стили компонента */
</style>
