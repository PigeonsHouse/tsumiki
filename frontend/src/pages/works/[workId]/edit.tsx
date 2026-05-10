import { useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import {
  useGetWork,
  useEditWork,
  useUpdateWorkThumbnail,
  useUploadThumbnail,
} from "../../../api";
import NotFound from "../../../components/NotFound";

type FormValues = {
  title: string;
  visibility: "public" | "limited";
  description: string;
  thumbnail: FileList;
};

const EditWork = () => {
  const { workId: workIdRaw } = useParams();
  const workID = useMemo(() => Number(workIdRaw), [workIdRaw]);
  const isValidId = !Number.isNaN(workID) && workID > 0;

  const navigate = useNavigate();
  const { data: work, isError } = useGetWork(workID, isValidId);
  const { mutateAsync: editWork } = useEditWork(workID);
  const { mutateAsync: uploadThumbnail } = useUploadThumbnail();
  const { mutateAsync: updateThumbnail } = useUpdateWorkThumbnail(workID);
  const [thumbnailPreview, setThumbnailPreview] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>();
  const thumbnailReg = register("thumbnail");

  useEffect(() => {
    if (!work) return;
    reset({
      title: work.title,
      visibility: work.visibility,
      description: work.description ?? "",
    });
  }, [work, reset]);

  const onSubmit = async (data: FormValues) => {
    const tasks: Promise<unknown>[] = [
      editWork({
        title: data.title,
        visibility: data.visibility,
        description: data.description || null,
      }),
    ];

    if (data.thumbnail?.length > 0) {
      tasks.push(
        uploadThumbnail(data.thumbnail[0]).then(({ id: thumbnailId }) =>
          updateThumbnail({ thumbnailId })
        )
      );
    }

    await Promise.all(tasks);
    navigate(`/works/${workID}`);
  };

  return !isValidId || isError ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>作品編集</h1>
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
          <label htmlFor="description">説明（4000文字以内）</label>
          <br />
          <textarea
            id="description"
            rows={5}
            style={{ width: "100%" }}
            {...register("description", {
              maxLength: {
                value: 4000,
                message: "4000文字以内で入力してください",
              },
            })}
          />
          {errors.description && (
            <p style={{ color: "red", margin: "4px 0 0" }}>
              {errors.description.message}
            </p>
          )}
        </div>

        <div style={{ marginBottom: 12 }}>
          <label htmlFor="thumbnail">
            サムネイル画像を差し替える（JPEG・PNG・GIF・最大10MB）
          </label>
          <br />
          {work?.thumbnailUrl && (
            <img
              src={work.thumbnailUrl}
              style={{
                width: "100%",
                aspectRatio: "16 / 9",
                objectFit: "contain",
                backgroundColor: "gray",
                marginBottom: 4,
              }}
            />
          )}
          <input
            id="thumbnail"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            {...thumbnailReg}
            onChange={(e) => {
              thumbnailReg.onChange(e);
              const file = e.target.files?.[0];
              setThumbnailPreview(file ? URL.createObjectURL(file) : null);
            }}
          />
          {thumbnailPreview && (
            <img
              src={thumbnailPreview}
              style={{
                marginTop: 4,
                width: "50%",
                aspectRatio: "16/9",
                objectFit: "contain",
                borderRadius: 4,
                backgroundColor: "gray",
                display: "block",
              }}
            />
          )}
        </div>

        <button type="submit" disabled={isSubmitting || !work}>
          {isSubmitting ? "保存中..." : "保存する"}
        </button>
      </form>
    </div>
  );
};

export default EditWork;
