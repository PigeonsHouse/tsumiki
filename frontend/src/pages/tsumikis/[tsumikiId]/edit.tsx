import { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import {
  useGetTsumiki,
  useEditTsumiki,
  useGetWorks,
  useUpdateTsumikiThumbnail,
  useUploadThumbnail,
} from "../../../api";

type FormValues = {
  title: string;
  visibility: "public" | "limited";
  workId: string;
  thumbnail: FileList;
};

const EditTsumiki = () => {
  const { tsumikiId: tsumikiIdRaw } = useParams();
  const tsumikiID = useMemo(() => Number(tsumikiIdRaw), [tsumikiIdRaw]);
  const isValidId = !Number.isNaN(tsumikiID) && tsumikiID > 0;

  const navigate = useNavigate();
  const { data: tsumiki } = useGetTsumiki(tsumikiID, isValidId);
  const { data: works } = useGetWorks();
  const { mutateAsync: editTsumiki } = useEditTsumiki(tsumikiID);
  const { mutateAsync: uploadThumbnail } = useUploadThumbnail();
  const { mutateAsync: updateThumbnail } = useUpdateTsumikiThumbnail(tsumikiID);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>();

  useEffect(() => {
    if (!tsumiki) return;
    reset({
      title: tsumiki.title,
      visibility: tsumiki.visibility,
      workId: tsumiki.work ? String(tsumiki.work.id) : "",
    });
  }, [tsumiki, reset]);

  const onSubmit = async (data: FormValues) => {
    const tasks: Promise<unknown>[] = [
      editTsumiki({
        title: data.title,
        visibility: data.visibility,
        workId: data.workId ? Number(data.workId) : null,
      }),
    ];

    if (data.thumbnail?.length > 0) {
      tasks.push(
        uploadThumbnail(data.thumbnail[0]).then(({ id: thumbnailId }) =>
          updateThumbnail({ thumbnailId }),
        ),
      );
    }

    await Promise.all(tasks);
    navigate(`/tsumikis/${tsumikiID}`);
  };

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>積み木編集</h1>
      <form
        onSubmit={handleSubmit(onSubmit)}
        style={{ maxWidth: 480, margin: "0 auto" }}
      >
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
          <label htmlFor="workId">紐づける作品</label>
          <br />
          <select
            id="workId"
            style={{ width: "100%" }}
            {...register("workId")}
          >
            <option value="">なし</option>
            {works?.map((work) => (
              <option key={work.id} value={work.id}>
                {work.title}
              </option>
            ))}
          </select>
        </div>

        <div style={{ marginBottom: 12 }}>
          <label htmlFor="thumbnail">サムネイル画像を差し替える（JPEG・PNG・GIF・最大5MB）</label>
          <br />
          {tsumiki?.thumbnailUrl && (
            <img
              src={tsumiki.thumbnailUrl}
              style={{ width: "100%", aspectRatio: "16 / 9", objectFit: "contain", backgroundColor: "gray", marginBottom: 4 }}
            />
          )}
          <input
            id="thumbnail"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            {...register("thumbnail")}
          />
        </div>

        <button type="submit" disabled={isSubmitting || !tsumiki}>
          {isSubmitting ? "保存中..." : "保存する"}
        </button>
      </form>
    </div>
  );
};

export default EditTsumiki;
