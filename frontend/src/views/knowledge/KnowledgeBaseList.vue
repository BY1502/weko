<template>
  <div class="kb-list-container">
    <div class="header">
      <div class="header-title">
        <h2>{{ $t("knowledgeBase.title") }}</h2>
        <p class="header-subtitle">{{ $t("knowledgeList.subtitle") }}</p>
      </div>
    </div>
    <div class="header-divider"></div>

    <div v-if="hasUninitializedKbs" class="warning-banner">
      <t-icon name="info-circle" size="16px" />
      <span>{{ $t("knowledgeList.uninitializedBanner") }}</span>
    </div>

    <div v-if="uploadSummaries.length" class="upload-progress-panel">
      <div
        v-for="summary in uploadSummaries"
        :key="summary.kbId"
        class="upload-progress-item"
      >
        <div class="upload-progress-icon">
          <t-icon
            :name="
              summary.completed === summary.total
                ? 'check-circle-filled'
                : 'upload'
            "
            size="20px"
          />
        </div>
        <div class="upload-progress-content">
          <div class="progress-title">
            {{
              summary.completed === summary.total
                ? $t("knowledgeList.uploadProgress.completedTitle", {
                    name: summary.kbName,
                  })
                : $t("knowledgeList.uploadProgress.uploadingTitle", {
                    name: summary.kbName,
                  })
            }}
          </div>
          <div class="progress-subtitle">
            {{
              summary.completed === summary.total
                ? $t("knowledgeList.uploadProgress.completedDetail", {
                    total: summary.total,
                  })
                : $t("knowledgeList.uploadProgress.detail", {
                    completed: summary.completed,
                    total: summary.total,
                  })
            }}
          </div>
          <div class="progress-subtitle secondary">
            {{
              summary.completed === summary.total
                ? $t("knowledgeList.uploadProgress.refreshing")
                : $t("knowledgeList.uploadProgress.keepPageOpen")
            }}
          </div>
          <div v-if="summary.hasError" class="progress-subtitle error">
            {{ $t("knowledgeList.uploadProgress.errorTip") }}
          </div>
          <div class="progress-bar">
            <div
              class="progress-bar-inner"
              :style="{ width: summary.progress + '%' }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="kbs.length > 0" class="kb-card-wrap">
      <div
        v-for="(kb, index) in kbs"
        :key="kb.id"
        class="kb-card"
        :class="{
          uninitialized: !isInitialized(kb),
          'kb-type-document': (kb.type || 'document') === 'document',
          'kb-type-faq': kb.type === 'faq',
          'highlight-flash':
            highlightedKbId !== null && highlightedKbId === kb.id,
        }"
        :ref="el => { if (highlightedKbId !== null && highlightedKbId === kb.id && el) highlightedCardRef = el as HTMLElement }"
        @click="handleCardClick(kb)"
      >
        <div class="card-header">
          <span class="card-title" :title="kb.name">{{ kb.name }}</span>
          <t-popup
            v-model="kb.showMore"
            overlayClassName="card-more-popup"
            :on-visible-change="onVisibleChange"
            trigger="click"
            destroy-on-close
            placement="bottom-right"
          >
            <div
              variant="outline"
              class="more-wrap"
              @click.stop="openMore(index)"
              :class="{ 'active-more': currentMoreIndex === index }"
            >
              <img class="more-icon" src="@/assets/img/more.png" alt="" />
            </div>
            <template #content>
              <div class="popup-menu" @click.stop>
                <div class="popup-menu-item" @click.stop="handleSettings(kb)">
                  <t-icon class="menu-icon" name="setting" />
                  <span>{{ $t("knowledgeBase.settings") }}</span>
                </div>
                <div
                  class="popup-menu-item delete"
                  @click.stop="handleDelete(kb)"
                >
                  <t-icon class="menu-icon" name="delete" />
                  <span>{{ $t("common.delete") }}</span>
                </div>
              </div>
            </template>
          </t-popup>
        </div>

        <div class="card-content">
          <div class="card-description">
            {{ kb.description || $t("knowledgeBase.noDescription") }}
          </div>
        </div>

        <div class="card-bottom">
          <div class="bottom-left">
            <div class="feature-badges">
              <t-tooltip
                :content="
                  kb.type === 'faq'
                    ? $t('knowledgeEditor.basic.typeFAQ')
                    : $t('knowledgeEditor.basic.typeDocument')
                "
                placement="top"
              >
                <div
                  class="feature-badge"
                  :class="{
                    'type-document': (kb.type || 'document') === 'document',
                    'type-faq': kb.type === 'faq',
                  }"
                >
                  <t-icon
                    :name="kb.type === 'faq' ? 'chat-bubble-help' : 'folder'"
                    size="14px"
                  />
                  <span class="badge-count">{{
                    kb.type === "faq"
                      ? kb.chunk_count || 0
                      : kb.knowledge_count || 0
                  }}</span>
                  <t-icon
                    v-if="kb.isProcessing"
                    name="loading"
                    size="12px"
                    class="processing-icon"
                  />
                </div>
              </t-tooltip>
              <t-tooltip
                v-if="kb.extract_config?.enabled"
                :content="$t('knowledgeList.features.knowledgeGraph')"
                placement="top"
              >
                <div class="feature-badge kg">
                  <t-icon name="relation" size="14px" />
                </div>
              </t-tooltip>
              <t-tooltip
                v-if="
                  kb.vlm_config?.enabled ||
                  (kb.cos_config?.provider && kb.cos_config?.bucket_name)
                "
                :content="$t('knowledgeList.features.multimodal')"
                placement="top"
              >
                <div class="feature-badge multimodal">
                  <t-icon name="image" size="14px" />
                </div>
              </t-tooltip>
              <t-tooltip
                v-if="kb.question_generation_config?.enabled"
                :content="$t('knowledgeList.features.questionGeneration')"
                placement="top"
              >
                <div class="feature-badge question">
                  <t-icon name="help-circle" size="14px" />
                </div>
              </t-tooltip>
            </div>
          </div>
          <span class="card-time">{{ kb.updated_at }}</span>
        </div>
      </div>
    </div>

    <div v-else-if="!loading" class="empty-state">
      <img class="empty-img" src="@/assets/img/upload.svg" alt="" />
      <span class="empty-txt">{{ $t("knowledgeList.empty.title") }}</span>
      <span class="empty-desc">{{
        $t("knowledgeList.empty.description")
      }}</span>
    </div>

    <t-dialog
      v-model:visible="deleteVisible"
      dialogClassName="del-knowledge-dialog"
      :closeBtn="false"
      :cancelBtn="null"
      :confirmBtn="null"
    >
      <div class="circle-wrap">
        <div class="dialog-header">
          <img class="circle-img" src="@/assets/img/circle.png" alt="" />
          <span class="circle-title">{{
            $t("knowledgeList.delete.confirmTitle")
          }}</span>
        </div>
        <span class="del-circle-txt">
          {{
            $t("knowledgeList.delete.confirmMessage", {
              name: deletingKb?.name ?? "",
            })
          }}
        </span>
        <div class="circle-btn">
          <span class="circle-btn-txt" @click="deleteVisible = false">{{
            $t("common.cancel")
          }}</span>
          <span class="circle-btn-txt confirm" @click="confirmDelete">{{
            $t("knowledgeList.delete.confirmButton")
          }}</span>
        </div>
      </div>
    </t-dialog>

    <KnowledgeBaseEditorModal
      :visible="uiStore.showKBEditorModal"
      :mode="uiStore.kbEditorMode"
      :kb-id="uiStore.currentKBId || undefined"
      :initial-type="uiStore.kbEditorType"
      @update:visible="(val) => (val ? null : uiStore.closeKBEditor())"
      @success="handleKBEditorSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed, watch, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { MessagePlugin, Icon as TIcon } from "tdesign-vue-next";
import { listKnowledgeBases, deleteKnowledgeBase } from "@/api/knowledge-base";
import { formatStringDate } from "@/utils/index";
import { useUIStore } from "@/stores/ui";
import KnowledgeBaseEditorModal from "./KnowledgeBaseEditorModal.vue";
import { useI18n } from "vue-i18n";

const router = useRouter();
const route = useRoute();
const uiStore = useUIStore();
const { t } = useI18n();

interface KB {
  id: string;
  name: string;
  description?: string;
  updated_at?: string;
  embedding_model_id?: string;
  summary_model_id?: string;
  type?: "document" | "faq";
  showMore?: boolean;
  vlm_config?: { enabled?: boolean; model_id?: string };
  extract_config?: { enabled?: boolean };
  cos_config?: { provider?: string; bucket_name?: string };
  question_generation_config?: { enabled?: boolean; question_count?: number };
  knowledge_count?: number;
  chunk_count?: number;
  isProcessing?: boolean; // 가져오기 작업 진행 중 여부
  processing_count?: number; // 처리 중인 문서 수 (문서 유형 전용)
}

const kbs = ref<KB[]>([]);
const loading = ref(false);
const deleteVisible = ref(false);
const deletingKb = ref<KB | null>(null);
const currentMoreIndex = ref<number>(-1);
const highlightedKbId = ref<string | null>(null);
const highlightedCardRef = ref<HTMLElement | null>(null);
const uploadTasks = ref<UploadTaskState[]>([]);
const uploadCleanupTimers = new Map<string, ReturnType<typeof setTimeout>>();
let uploadRefreshTimer: ReturnType<typeof setTimeout> | null = null;
const UPLOAD_CLEANUP_DELAY = 10000;

interface UploadTaskState {
  uploadId: string;
  kbId: string;
  fileName?: string;
  progress: number;
  status: "uploading" | "success" | "error";
  error?: string;
}

interface UploadSummary {
  kbId: string;
  kbName: string;
  total: number;
  completed: number;
  progress: number;
  hasError: boolean;
}

const fetchList = () => {
  loading.value = true;
  return listKnowledgeBases()
    .then((res: any) => {
      const data = res.data || [];
      // 시간 포맷팅 및 showMore 상태 초기화
      // is_processing 필드는 백엔드에서 반환됨
      kbs.value = data.map((kb: any) => ({
        ...kb,
        updated_at: kb.updated_at
          ? formatStringDate(new Date(kb.updated_at))
          : "",
        showMore: false,
        isProcessing: kb.is_processing || false,
        processing_count: kb.processing_count || 0,
      }));
    })
    .finally(() => (loading.value = false));
};

onMounted(() => {
  fetchList().then(() => {
    // 라우트 파라미터에서 하이라이트할 지식 베이스 ID 확인
    const highlightKbId = route.query.highlightKbId as string;
    if (highlightKbId) {
      triggerHighlightFlash(highlightKbId);
      // URL 쿼리 파라미터 제거
      router.replace({ query: {} });
    }
  });

  window.addEventListener(
    "knowledgeFileUploadStart",
    handleUploadStartEvent as EventListener
  );
  window.addEventListener(
    "knowledgeFileUploadProgress",
    handleUploadProgressEvent as EventListener
  );
  window.addEventListener(
    "knowledgeFileUploadComplete",
    handleUploadCompleteEvent as EventListener
  );
  window.addEventListener(
    "knowledgeFileUploaded",
    handleUploadFinishedEvent as EventListener
  );
});

onUnmounted(() => {
  window.removeEventListener(
    "knowledgeFileUploadStart",
    handleUploadStartEvent as EventListener
  );
  window.removeEventListener(
    "knowledgeFileUploadProgress",
    handleUploadProgressEvent as EventListener
  );
  window.removeEventListener(
    "knowledgeFileUploadComplete",
    handleUploadCompleteEvent as EventListener
  );
  window.removeEventListener(
    "knowledgeFileUploaded",
    handleUploadFinishedEvent as EventListener
  );

  uploadCleanupTimers.forEach((timer) => clearTimeout(timer));
  uploadCleanupTimers.clear();
  if (uploadRefreshTimer) {
    clearTimeout(uploadRefreshTimer);
    uploadRefreshTimer = null;
  }
});

// 라우트 변경 감지, 다른 페이지에서 넘어온 하이라이트 요청 처리
watch(
  () => route.query.highlightKbId,
  (newKbId) => {
    if (newKbId && typeof newKbId === "string" && kbs.value.length > 0) {
      triggerHighlightFlash(newKbId);
      router.replace({ query: {} });
    }
  }
);

const openMore = (index: number) => {
  // 현재 열린 인덱스만 기록하여 활성화 스타일 표시
  // 팝업 창의 열림/닫힘은 v-model로 자동 관리됨
  currentMoreIndex.value = index;
};

const onVisibleChange = (visible: boolean) => {
  // 팝업 창이 닫힐 때 인덱스 초기화
  if (!visible) {
    currentMoreIndex.value = -1;
  }
};

const handleSettings = (kb: KB) => {
  // 수동으로 팝업 창 닫기
  kb.showMore = false;
  goSettings(kb.id);
};

const handleDelete = (kb: KB) => {
  // 수동으로 팝업 창 닫기
  kb.showMore = false;
  deletingKb.value = kb;
  deleteVisible.value = true;
};

const confirmDelete = () => {
  if (!deletingKb.value) return;

  deleteKnowledgeBase(deletingKb.value.id)
    .then((res: any) => {
      if (res.success) {
        MessagePlugin.success(t("knowledgeList.messages.deleted"));
        deleteVisible.value = false;
        deletingKb.value = null;
        fetchList();
      } else {
        MessagePlugin.error(
          res.message || t("knowledgeList.messages.deleteFailed")
        );
      }
    })
    .catch((e: any) => {
      MessagePlugin.error(
        e?.message || t("knowledgeList.messages.deleteFailed")
      );
    });
};

const isInitialized = (kb: KB) => {
  return !!(
    kb.embedding_model_id &&
    kb.embedding_model_id !== "" &&
    kb.summary_model_id &&
    kb.summary_model_id !== ""
  );
};

// 초기화되지 않은 지식 베이스가 있는지 계산
const hasUninitializedKbs = computed(() => {
  return kbs.value.some((kb) => !isInitialized(kb));
});

const getKbDisplayName = (kbId: string) => {
  const target = kbs.value.find((kb) => kb.id === kbId);
  if (target?.name) return target.name;
  return t("knowledgeList.uploadProgress.unknownKb", { id: kbId }) as string;
};

const uploadSummaries = computed<UploadSummary[]>(() => {
  if (!uploadTasks.value.length) return [];
  const grouped: Record<string, UploadTaskState[]> = {};
  uploadTasks.value.forEach((task) => {
    const kbKey = String(task.kbId);
    if (!grouped[kbKey]) grouped[kbKey] = [];
    grouped[kbKey].push(task);
  });
  return Object.entries(grouped)
    .map(([kbId, tasks]) => {
      const total = tasks.length;
      const completed = tasks.filter(
        (task) => task.status !== "uploading"
      ).length;
      const progressSum = tasks.reduce(
        (sum, task) => sum + (task.progress ?? 0),
        0
      );
      const avgProgress =
        total === 0
          ? 0
          : Math.min(100, Math.max(0, Math.round(progressSum / total)));
      const hasError = tasks.some((task) => task.status === "error");
      return {
        kbId,
        kbName: getKbDisplayName(kbId),
        total,
        completed,
        progress: avgProgress,
        hasError,
      };
    })
    .sort((a, b) => a.kbName.localeCompare(b.kbName));
});

const clampProgress = (value: number) =>
  Math.min(100, Math.max(0, Math.round(value)));

const addUploadTask = (task: UploadTaskState) => {
  uploadTasks.value = [
    ...uploadTasks.value.filter((item) => item.uploadId !== task.uploadId),
    task,
  ];
};

const patchUploadTask = (uploadId: string, patch: Partial<UploadTaskState>) => {
  const index = uploadTasks.value.findIndex(
    (task) => task.uploadId === uploadId
  );
  if (index === -1) return;
  const nextTasks = [...uploadTasks.value];
  nextTasks[index] = { ...nextTasks[index], ...patch };
  uploadTasks.value = nextTasks;
};

const removeUploadTask = (uploadId: string) => {
  uploadTasks.value = uploadTasks.value.filter(
    (task) => task.uploadId !== uploadId
  );
  const timer = uploadCleanupTimers.get(uploadId);
  if (timer) {
    clearTimeout(timer);
    uploadCleanupTimers.delete(uploadId);
  }
};

const scheduleUploadTaskCleanup = (uploadId: string) => {
  const existing = uploadCleanupTimers.get(uploadId);
  if (existing) {
    clearTimeout(existing);
  }
  const timer = setTimeout(() => {
    removeUploadTask(uploadId);
  }, UPLOAD_CLEANUP_DELAY);
  uploadCleanupTimers.set(uploadId, timer);
};

type UploadEventDetail = {
  uploadId: string;
  kbId?: string | number;
  fileName?: string;
  progress?: number;
  status?: UploadTaskState["status"];
  error?: string;
};

const ensureUploadTaskEntry = (detail?: UploadEventDetail) => {
  if (!detail?.uploadId) return null;
  const existing = uploadTasks.value.find(
    (task) => task.uploadId === detail.uploadId
  );
  if (existing) return existing;
  if (!detail.kbId) return null;
  const initialProgress =
    typeof detail.progress === "number" ? clampProgress(detail.progress) : 0;
  const newTask: UploadTaskState = {
    uploadId: detail.uploadId,
    kbId: String(detail.kbId),
    fileName: detail.fileName,
    progress: initialProgress,
    status: detail.status || "uploading",
    error: detail.error,
  };
  addUploadTask(newTask);
  return newTask;
};

const handleCardClick = (kb: KB) => {
  if (isInitialized(kb)) {
    goDetail(kb.id);
  } else {
    goSettings(kb.id);
  }
};

const goDetail = (id: string) => {
  router.push(`/platform/knowledge-bases/${id}`);
};

const goSettings = (id: string) => {
  // 모달 창을 사용하여 설정 열기
  uiStore.openKBSettings(id);
};

// 지식 베이스 에디터 성공 콜백 (생성 또는 편집 성공)
const handleKBEditorSuccess = (kbId: string) => {
  console.log("[KnowledgeBaseList] knowledge operation success:", kbId);
  fetchList().then(() => {
    // 라우트 파라미터에서 가져온 하이라이트 ID인 경우, 깜빡임 효과 트리거
    if (route.query.highlightKbId === kbId) {
      triggerHighlightFlash(kbId);
      // URL 쿼리 파라미터 제거
      router.replace({ query: {} });
    }
  });
};

// 하이라이트 깜빡임 효과 트리거
const triggerHighlightFlash = (kbId: string) => {
  highlightedKbId.value = kbId;
  nextTick(() => {
    if (highlightedCardRef.value) {
      // 하이라이트된 카드로 스크롤
      highlightedCardRef.value.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    }
    // 3초 후 하이라이트 제거
    setTimeout(() => {
      highlightedKbId.value = null;
    }, 3000);
  });
};

const handleUploadStartEvent = (event: Event) => {
  const detail = (event as CustomEvent<UploadEventDetail>).detail;
  if (!detail?.uploadId || !detail?.kbId) return;
  addUploadTask({
    uploadId: detail.uploadId,
    kbId: String(detail.kbId),
    fileName: detail.fileName,
    progress:
      typeof detail.progress === "number" ? clampProgress(detail.progress) : 0,
    status: "uploading",
  });
};

const handleUploadProgressEvent = (event: Event) => {
  const detail = (event as CustomEvent<UploadEventDetail>).detail;
  if (!detail?.uploadId || typeof detail.progress !== "number") return;
  if (!ensureUploadTaskEntry(detail)) return;
  patchUploadTask(detail.uploadId, {
    progress: clampProgress(detail.progress),
  });
};

const handleUploadCompleteEvent = (event: Event) => {
  const detail = (event as CustomEvent<UploadEventDetail>).detail;
  if (!detail?.uploadId) return;
  const progress =
    typeof detail.progress === "number" ? clampProgress(detail.progress) : 100;
  if (!ensureUploadTaskEntry({ ...detail, progress })) return;
  patchUploadTask(detail.uploadId, {
    status: detail.status || "success",
    progress,
    error: detail.error,
  });
  scheduleUploadTaskCleanup(detail.uploadId);
};

const handleUploadFinishedEvent = (event: Event) => {
  const detail = (event as CustomEvent<{ kbId?: string | number }>).detail;
  if (!detail?.kbId) return;
  if (uploadRefreshTimer) {
    clearTimeout(uploadRefreshTimer);
  }
  uploadRefreshTimer = setTimeout(() => {
    fetchList();
    uploadRefreshTimer = null;
  }, 800);
};
</script>

<style scoped lang="less">
/* ... (스타일 섹션은 변경 사항 없음) ... */
.kb-list-container {
  padding: 24px 44px;
  // background: #fff;
  margin: 0 20px;
  height: calc(100vh);
  overflow-y: auto;
  box-sizing: border-box;
  flex: 1;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;

  .header-title {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  h2 {
    margin: 0;
    color: #000000e6;
    font-family: "PingFang SC";
    font-size: 24px;
    font-weight: 600;
    line-height: 32px;
  }
}

.header-subtitle {
  margin: 0;
  color: #00000099;
  font-family: "PingFang SC";
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
}

.header-divider {
  height: 1px;
  background: #e7ebf0;
  margin-bottom: 20px;
}

/* ... (이하 스타일 생략) ... */

// 반응형 레이아웃
@media (min-width: 900px) {
  .kb-card-wrap {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1250px) {
  .kb-card-wrap {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1600px) {
  .kb-card-wrap {
    grid-template-columns: repeat(4, 1fr);
  }
}

// 삭제 확인 대화상자 스타일
:deep(.del-knowledge-dialog) {
  padding: 0px !important;
  border-radius: 6px !important;

  .t-dialog__header {
    display: none;
  }

  .t-dialog__body {
    padding: 16px;
  }

  .t-dialog__footer {
    padding: 0;
  }
}

/* ... */
</style>

<style lang="less">
// 더보기 작업 팝업 스타일
.card-more-popup {
  z-index: 99 !important;

  .t-popup__content {
    padding: 6px 0 !important;
    margin-top: 6px !important;
    min-width: 140px;
    border-radius: 6px !important;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1) !important;
    border: 1px solid #e7ebf0 !important;
  }
}

/* ... */

// 생성 대화상자 스타일 최적화
.create-kb-dialog {
  .t-form-item__label {
    font-family: "PingFang SC";
    font-size: 14px;
    font-weight: 500;
    color: #000000e6;
  }
  /* ... */
}
</style>
