import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { useCreateWork } from "../../api";
import { thumbnailsApi } from "../../api/client";

type FormValues = {
  title: string;
  visibility: "public" | "limited";
  description: string;
  thumbnail: FileList;
};

const NewWork = () => {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    defaultValues: { visibility: "public" },
  });

  const { mutateAsync: createWork } = useCreateWork();

  const onSubmit = async (data: FormValues) => {
    let thumbnailId: number | null = null;
    if (data.thumbnail?.length > 0) {
      const { id } = await thumbnailsApi.postThumbnail({
        thumbnail: data.thumbnail[0],
      });
      thumbnailId = id;
    }

    const work = await createWork({
      title: data.title,
      visibility: data.visibility,
      description: data.description || undefined,
      thumbnailId,
    });

    navigate(`/works/${work.id}`);
  };

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>作品新規作成</h1>
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
          <label htmlFor="description">説明（4000文字以内）</label>
          <br />
          <textarea
            id="description"
            rows={5}
            style={{ width: "100%" }}
            {...register("description", {
              maxLength: { value: 4000, message: "4000文字以内で入力してください" },
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
            サムネイル画像（JPEG・PNG・GIF・最大5MB）
          </label>
          <br />
          <input
            id="thumbnail"
            type="file"
            accept="image/jpeg,image/png,image/gif"
            {...register("thumbnail")}
          />
        </div>

        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "作成中..." : "作成する"}
        </button>
      </form>
    </div>
  );
};

export default NewWork;
