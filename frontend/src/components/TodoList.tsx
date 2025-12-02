import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import TodoItem from "./TodoItem";
import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";

// const BASE_URL = "http://localhost:5000/api"; // آدرس سرور

const BASE_URL = `${import.meta.env.VITE_API_URL || 'http://localhost:5000'}/api`;

type Todo = {
	_id: string;
	body: string;
	complete: boolean;
};

const TodoList = () => {
	const queryClient = useQueryClient();

	// ✅ گرفتن لیست تودو‌ها
	const { data: todos = [], isLoading } = useQuery<Todo[]>({
		queryKey: ["todos"],
		queryFn: async () => {
			console.log("👉 Fetching todos...");
			const res = await fetch(BASE_URL + "/todos");
			const data = await res.json();
			console.log("✅ Todos fetched:", data);
			return data.map(todo => ({
				...todo,
				complete: todo.complete || false, // مقدار کامل رو تنظیم کن
			}));
		},
	});

	// ✅ حذف تودو
	const deleteTodo = async (id: string) => {
		await fetch(`${BASE_URL}/todos/${id}`, {
			method: "DELETE",
		});

		// ✅ آپدیت کش بعد از حذف
		queryClient.setQueryData<Todo[]>(["todos"], (oldTodos) => {
			return oldTodos?.filter(todo => todo._id !== id) || [];
		});
	};

	// ✅ آپدیت تودو
	const updateTodo = async (id: string) => {
		const res = await fetch(`${BASE_URL}/todos/${id}`, {
			method: "PATCH",
		});
		const updatedTodo = await res.json();

		// ✅ آپدیت کش بعد از آپدیت
		queryClient.setQueryData<Todo[]>(["todos"], (oldTodos) => {
			if (!oldTodos) return [];
			return oldTodos.map((todo) =>
				todo._id === updatedTodo._id ? { ...todo, ...updatedTodo } : todo
			);
		});
	};

	if (isLoading) return(
		<Flex justifyContent={"center"} my={4}>
			<Spinner size={"xl"} />
		</Flex>
	);

	return (
		<div className="p-4">
			<Text
				fontSize={"4xl"}
				textTransform={"uppercase"}
				fontWeight={"bold"}
				textAlign={"center"}
				my={2}
				bgGradient='linear(to-l, #0b85f8, #00ffff)'
				bgClip='text'
			>
				Today's Tasks
			</Text>
			{todos.length === 0 ? (
				// <div className="flex justify-center items-center">
				// 	<span className="text-xl text-gray-500">No tasks available</span>
				// </div>

				<Stack alignItems={"center"} gap='3'>
					<Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
						All tasks completed! 🤞
					</Text>
					<img src='/go.png' alt='Go logo' width={70} height={70} />
				</Stack>
			) : (
				todos.map((todo) => (
					<TodoItem
						key={todo._id}
						todo={todo}
						onDelete={deleteTodo}
						onUpdate={updateTodo}
					/>
				))
			)}
		</div>
	);
};

export default TodoList;
