import { ref } from 'vue';
import { models } from '../../wailsjs/go/models';
import { GetAllLoans, GetLoanByID, CreateLoan, UpdateLoan, DeleteLoan } from '../../wailsjs/go/main/App';

export function useLoans() {
  const loans = ref<models.Loan[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // 加载所有贷款
  const loadLoans = async () => {
    loading.value = true;
    error.value = null;
    try {
      const result = await GetAllLoans();
      loans.value = result || [];
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载贷款列表失败';
      console.error('加载贷款失败:', err);
    } finally {
      loading.value = false;
    }
  };

  // 根据ID获取贷款
  const getLoan = async (id: number): Promise<models.Loan | null> => {
    loading.value = true;
    error.value = null;
    try {
      const result = await GetLoanByID(id);
      return result;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载贷款详情失败';
      console.error('加载贷款详情失败:', err);
      return null;
    } finally {
      loading.value = false;
    }
  };

  // 创建贷款
  const createLoan = async (loan: Partial<models.Loan>): Promise<boolean> => {
    loading.value = true;
    error.value = null;
    try {
      await CreateLoan(loan as models.Loan);
      await loadLoans();
      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建贷款失败';
      console.error('创建贷款失败:', err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 更新贷款
  const updateLoan = async (loan: models.Loan): Promise<boolean> => {
    loading.value = true;
    error.value = null;
    try {
      await UpdateLoan(loan);
      await loadLoans();
      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新贷款失败';
      console.error('更新贷款失败:', err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 删除贷款
  const deleteLoan = async (id: number): Promise<boolean> => {
    loading.value = true;
    error.value = null;
    try {
      await DeleteLoan(id);
      await loadLoans();
      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '删除贷款失败';
      console.error('删除贷款失败:', err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  return {
    loans,
    loading,
    error,
    loadLoans,
    getLoan,
    createLoan,
    updateLoan,
    deleteLoan,
  };
}
