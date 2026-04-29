import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { useUploadThumbnail, useCreateTsumiki } from "../../api";
import { blocksApi, tsumikisApi } from "../../api/client";

type FormValues = {
  title: string;
  visibility: "public" | "limited";
  thumbnail: FileList;
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
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ defaultValues: { visibility: "public", percentage: 0, condition: 3 } });

  const { mutateAsync: uploadThumbnail } = useUploadThumbnail();
  const { mutateAsync: createTsumiki } = useCreateTsumiki();

  const onSubmit = async (data: FormValues) => {
    // 1. サムネイルアップロード
    const { id: thumbnailId } = await uploadThumbnail(data.thumbnail[0]);

    // 2. 積み木作成
    const tsumiki = await createTsumiki({
      title: data.title,
      visibility: data.visibility,
      thumbnailId,
    });
    const tsumikiID = tsumiki.id;

    // 3. メディアアップロード（選択されていれば最大4件）
    const mediaIds: number[] = [];
    if (data.medias && data.medias.length > 0) {
      const files = Array.from(data.medias).slice(0, 4);
      for (const file of files) {
        const media = await tsumikisApi.postMedia({ tsumikiID, file });
        mediaIds.push(media.id);
      }
    }

    // 4. 最初のブロック追加
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
      <form onSubmit={handleSubmit(onSubmit)} style={{ maxWidth: 480, margin: "0 auto" }}>
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
                maxLength: { value: 200, message: "200文字以内で入力してください" },
              })}
            />
            {errors.title && <p style={{ color: "red", margin: "4px 0 0" }}>{errors.title.message}</p>}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="visibility">公開設定（必須）</label>
            <br />
            <select id="visibility" style={{ width: "100%" }} {...register("visibility")}>
              <option value="public">public（誰でも閲覧可能）</option>
              <option value="limited">limited（同サーバーメンバーのみ）</option>
            </select>
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="thumbnail">サムネイル画像（必須 / JPEG・PNG・GIF・最大5MB）</label>
            <br />
            <input
              id="thumbnail"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              {...register("thumbnail", { required: "サムネイルは必須です" })}
            />
            {errors.thumbnail && <p style={{ color: "red", margin: "4px 0 0" }}>{errors.thumbnail.message}</p>}
          </div>
        </fieldset>

        <fieldset style={{ marginBottom: 24 }}>
          <legend>最初のブロック</legend>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="percentage">進捗率（必須・0〜100）</label>
            <br />
            <input
              id="percentage"
              type="number"
              min={0}
              max={100}
              style={{ width: "100%" }}
              {...register("percentage", {
                required: "進捗率は必須です",
                min: { value: 0, message: "0以上で入力してください" },
                max: { value: 100, message: "100以下で入力してください" },
              })}
            />
            {errors.percentage && <p style={{ color: "red", margin: "4px 0 0" }}>{errors.percentage.message}</p>}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="condition">コンディション（必須・1〜5）</label>
            <br />
            <select id="condition" style={{ width: "100%" }} {...register("condition", { required: true })}>
              <option value={1}>1</option>
              <option value={2}>2</option>
              <option value={3}>3</option>
              <option value={4}>4</option>
              <option value={5}>5</option>
            </select>
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="message">メッセージ（任意・200文字以内）</label>
            <br />
            <textarea
              id="message"
              rows={3}
              style={{ width: "100%" }}
              {...register("message", {
                maxLength: { value: 200, message: "200文字以内で入力してください" },
                validate: (value) => {
                  const medias = getValues("medias");
                  return !!value || medias?.length > 0 || "メッセージかメディアのいずれかは必須です";
                },
              })}
            />
            {errors.message && <p style={{ color: "red", margin: "4px 0 0" }}>{errors.message.message}</p>}
          </div>

          <div style={{ marginBottom: 12 }}>
            <label htmlFor="medias">メディア（任意・最大4件 / 画像・音声・動画）</label>
            <br />
            <input
              id="medias"
              type="file"
              multiple
              accept="image/jpeg,image/png,image/gif,audio/mpeg,audio/wav,audio/ogg,audio/aac,video/mp4,video/webm,video/quicktime"
              {...register("medias", {
                validate: (value) => {
                  const message = getValues("message");
                  return !!message || value?.length > 0 || "メッセージかメディアのいずれかは必須です";
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
