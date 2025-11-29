import { useState } from "react";
import { FaCheckCircle, FaRegCircle, FaTrash } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";

type Todo = {
	_id: string;
	body: string;
	complete: boolean;
};

type Props = {
	todo: Todo;
	onDelete: (id: string) => void;
	onUpdate: (id: string) => void;
};

const TodoItem = ({ todo, onDelete, onUpdate }: Props) => {
	const [completed, setCompleted] = useState(todo.complete);

	const handleUpdate = async () => {
		console.log("👉 Update clicked");
		await onUpdate(todo._id);
		setCompleted((prev) => !prev);
	};

	return (
		<Flex gap={2} alignItems={"center"}>
			<Flex
				flex={1}
				alignItems={"center"}
				border={"1px"}
				borderColor={"gray.600"}
				p={2}
				borderRadius={"lg"}
				justifyContent={"space-between"}
			>
				<Text
					color={completed ? "green.200" : "yellow.100"}
					textDecoration={completed ? "line-through" : "none"}
				>
					{todo.body}
				</Text>
				{completed && (
					<Badge ml='1' colorScheme='green'>
						Done
					</Badge>
				)}
				{!completed && (
					<Badge ml='1' colorScheme='yellow'>
						In Progress
					</Badge>
				)}
			</Flex>
			<Flex gap={2} alignItems={"center"}>
				<Box color={"green.500"} cursor={"pointer"} onClick={handleUpdate}>
					{completed && <FaCheckCircle size={20} />}
					{!completed && <FaRegCircle size={20}/>}
				</Box>
				<Box color={"red.500"} cursor={"pointer"} onClick={() => onDelete(todo._id)}>
					<MdDelete size={25} />
				</Box>
			</Flex>
		</Flex>
	);
};

export default TodoItem;
