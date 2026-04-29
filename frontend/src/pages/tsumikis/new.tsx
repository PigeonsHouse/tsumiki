import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import {
  useUploadThumbnail,
  useCreateTsumiki,
  useGetWorks,
  useCreateWork,
} from "../../api";
import { blocksApi, tsumikisApi } from "../../api/client";
import { thumbnailsApi } from "../../api/client";

type FormValues = {
  title: string;
  visibility: "public" | "limited";
  thumbnail: FileList;
  workMode: "none" | "existing" | "new";
  workId: string;
  newWorkTitle: string;
  newWorkVisibility: "public" | "limited";
  newWorkDescription: string;
  newWorkThumbnail: FileList;
  message: string;
  percentage: number;
  condition: number;
  medias: FileList;
};

const NewTsumiki = () => {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    watch,
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    defaultValues: {
      visibility: "public",
      workMode: "none",
      newWorkVisibility: "public",
      percentage: 0,
      condition: 3,
    },
  });

  const workMode = watch("workMode");
  const percentage = watch("percentage");
  const condition = watch("condition");
  const { data: works } = useGetWorks();
  const { mutateAsync: uploadThumbnail } = useUploadThumbnail();
  const { mutateAsync: createTsumiki } = useCreateTsumiki();
  const { mutateAsync: createWork } = useCreateWork();

  const onSubmit = async (data: FormValues) => {
    // 1. サムネイルアップロード
    const { id: thumbnailId } = await uploadThumbnail(data.thumbnail[0]);

    // 2. 作品の解決（既存選択 or 新規作成）
    let workId: number | null = null;
    if (data.workMode === "existing" && data.workId) {
      workId = Number(data.workId);
    } else if (data.workMode === "new") {
      let newWorkThumbnailId: number | null = null;
      if (data.newWorkThumbnail?.length > 0) {
        const { id } = await thumbnailsApi.postThumbnail({
          thumbnail: data.newWorkThumbnail[0],
        });
        newWorkThumbnailId = id;
      }
      const work = await createWork({
        title: data.newWorkTitle,
        visibility: data.newWorkVisibility,
        description: data.newWorkDescription,
        thumbnailId: newWorkThumbnailId,
      });
      workId = work.id;
    }

    // 3. 積み木作成
    const tsumiki = await createTsumiki({
      title: data.title,
      visibility: data.visibility,
      thumbnailId,
      workId,
    });
    const tsumikiID = tsumiki.id;

    // 4. メディアアップロード（選択されていれば最大4件）
    const mediaIds: number[] = [];
    if (data.medias && data.medias.length > 0) {
      const files = Array.from(data.medias).slice(0, 4);
      for (const file of files) {
        const media = await tsumikisApi.postMedia({ tsumikiID, file });
        mediaIds.push(media.id);
      }
    }

    // 5. 最初のブロック追加
    await blocksApi.addBlock({
      tsumikiID,
      addBlockRequest: {
        message: data.message || null,
        percentage: Number(data.percentage),
        condition: Number(data.condition),
        mediaIds: mediaIds.length > 0 ? mediaIds : undefined,
        latestBlockId: null,
      },
    });

    navigate(`/tsumikis/${tsumikiID}`);
  };

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>積み木新規作成</h1>
      <form
        onSubmit={handleSubmit(onSubmit)}
        style={{ maxWidth: 480, margin: "0 auto" }}
      >
        <fieldset style={{ marginBottom: 24 }}>
          <legend>積み木情報</legend>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="title">タイトル（必須）</label>
            <br />
            <input
              id="title"
              type="text"
              style={{ width: "100%" }}
              {...register("title", {
                required: "タイトルは必須です",
                maxLength: {
                  value: 200,
                  message: "200文字以内で入力してください",
                },
              })}
            />
            {errors.title && (
              <p style={{ color: "red", margin: "4px 0 0" }}>
                {errors.title.message}
              </p>
            )}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="visibility">公開設定（必須）</label>
            <br />
            <select
              id="visibility"
              style={{ width: "100%" }}
              {...register("visibility")}
            >
              <option value="public">public（誰でも閲覧可能）</option>
              <option value="limited">limited（同サーバーメンバーのみ）</option>
            </select>
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="thumbnail">
              サムネイル画像（必須 / JPEG・PNG・GIF・最大5MB）
            </label>
            <br />
            <input
              id="thumbnail"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              {...register("thumbnail", { required: "サムネイルは必須です" })}
            />
            {errors.thumbnail && (
              <p style={{ color: "red", margin: "4px 0 0" }}>
                {errors.thumbnail.message}
              </p>
            )}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label>紐づける作品（任意）</label>
            <br />
            <label>
              <input type="radio" value="none" {...register("workMode")} /> なし
            </label>{" "}
            <label>
              <input type="radio" value="existing" {...register("workMode")} />{" "}
              既存から選ぶ
            </label>{" "}
            <label>
              <input type="radio" value="new" {...register("workMode")} />{" "}
              新しく作成
            </label>
          </div>

          {workMode === "existing" && (
            <div style={{ marginBottom: 12 }}>
              <label htmlFor="workId">作品を選択</label>
              <br />
              <select
                id="workId"
                style={{ width: "100%" }}
                {...register("workId")}
              >
                <option value="">選択してください</option>
                {works?.map((work) => (
                  <option key={work.id} value={work.id}>
                    {work.title}
                  </option>
                ))}
              </select>
            </div>
          )}

          {workMode === "new" && (
            <fieldset style={{ marginBottom: 12 }}>
              <legend>新しい作品</legend>

              <div style={{ marginBottom: 8 }}>
                <label htmlFor="newWorkTitle">タイトル（必須）</label>
                <br />
                <input
                  id="newWorkTitle"
                  type="text"
                  style={{ width: "100%" }}
                  {...register("newWorkTitle", {
                    required:
                      workMode === "new" ? "作品タイトルは必須です" : false,
                    maxLength: {
                      value: 200,
                      message: "200文字以内で入力してください",
                    },
                  })}
                />
                {errors.newWorkTitle && (
                  <p style={{ color: "red", margin: "4px 0 0" }}>
                    {errors.newWorkTitle.message}
                  </p>
                )}
              </div>

              <div style={{ marginBottom: 8 }}>
                <label htmlFor="newWorkVisibility">公開設定（必須）</label>
                <br />
                <select
                  id="newWorkVisibility"
                  style={{ width: "100%" }}
                  {...register("newWorkVisibility")}
                >
                  <option value="public">public（誰でも閲覧可能）</option>
                  <option value="limited">
                    limited（同サーバーメンバーのみ）
                  </option>
                </select>
              </div>

              <div style={{ marginBottom: 8 }}>
                <label htmlFor="newWorkDescription">説明（必須）</label>
                <br />
                <textarea
                  id="newWorkDescription"
                  rows={3}
                  style={{ width: "100%" }}
                  {...register("newWorkDescription", {
                    required: workMode === "new" ? "説明は必須です" : false,
                    maxLength: {
                      value: 4000,
                      message: "4000文字以内で入力してください",
                    },
                  })}
                />
                {errors.newWorkDescription && (
                  <p style={{ color: "red", margin: "4px 0 0" }}>
                    {errors.newWorkDescription.message}
                  </p>
                )}
              </div>

              <div style={{ marginBottom: 8 }}>
                <label htmlFor="newWorkThumbnail">サムネイル画像（任意）</label>
                <br />
                <input
                  id="newWorkThumbnail"
                  type="file"
                  accept="image/jpeg,image/png,image/gif"
                  {...register("newWorkThumbnail")}
                />
              </div>
            </fieldset>
          )}
        </fieldset>

        <fieldset style={{ marginBottom: 24 }}>
          <legend>最初のブロック</legend>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="percentage">進捗率: {percentage}%</label>
            <br />
            <input
              id="percentage"
              type="range"
              min={0}
              max={100}
              step={1}
              style={{ width: "100%" }}
              {...register("percentage", { required: true })}
            />
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="condition">コンディション: {condition} / 5</label>
            <br />
            <input
              id="condition"
              type="range"
              min={1}
              max={5}
              step={1}
              style={{ width: "100%" }}
              {...register("condition", { required: true })}
            />
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="message">メッセージ（任意・200文字以内）</label>
            <br />
            <textarea
              id="message"
              rows={3}
              style={{ width: "100%" }}
              {...register("message", {
                maxLength: {
                  value: 200,
                  message: "200文字以内で入力してください",
                },
                validate: (value) => {
                  const medias = getValues("medias");
                  return (
                    !!value ||
                    medias?.length > 0 ||
                    "メッセージかメディアのいずれかは必須です"
                  );
                },
              })}
            />
            {errors.message && (
              <p style={{ color: "red", margin: "4px 0 0" }}>
                {errors.message.message}
              </p>
            )}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="medias">
              メディア（任意・最大4件 / 画像・音声・動画）
            </label>
            <br />
            <input
              id="medias"
              type="file"
              multiple
              accept="image/jpeg,image/png,image/gif,audio/mpeg,audio/wav,audio/ogg,audio/aac,video/mp4,video/webm,video/quicktime"
              {...register("medias", {
                validate: (value) => {
                  const message = getValues("message");
                  return (
                    !!message ||
                    value?.length > 0 ||
                    "メッセージかメディアのいずれかは必須です"
                  );
                },
              })}
            />
          </div>
        </fieldset>

        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "作成中..." : "作成する"}
        </button>
      </form>
    </div>
  );
};

export default NewTsumiki;
