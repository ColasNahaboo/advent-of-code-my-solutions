#!/usr/bin/python3
# From deividragon
# https://github.com/deivi-drg/advent-of-code-2022/blob/main/Day22/day22.py
# https://www.reddit.com/r/adventofcode/comments/zsct8w/comment/j28g7ue/?utm_source=reddit&utm_medium=web2x&context=3

# Colas: I have modified it to emit traces in the exact same format as my code.
# I could thus compare the steps of it with my solution for debugging

import sys
# ignore broken pipe on |head
from signal import signal, SIGPIPE, SIG_DFL
signal(SIGPIPE, SIG_DFL)
    
DIRECTION_VALUES = [(1, 0), (0, 1), (-1, 0), (0, -1)]
# Right, down, left, up.
# Anlges are 0 for upright, then clockwise.


def parse_instructions(instructions: str) -> list:
    parsed_instructions = []
    number_buffer = ""
    for character in instructions:
        if character in ["R", "L"]:
            parsed_instructions.append(int(number_buffer))
            parsed_instructions.append(character)
            number_buffer = ""
        else:
            number_buffer += character
    if number_buffer != "":
        parsed_instructions.append(int(number_buffer))
    return parsed_instructions


def sum_tuples(tuple_1: tuple[int, int], tuple_2: tuple[int, int]) -> tuple[int, int]:
    return (tuple_1[0] + tuple_2[0], tuple_1[1] + tuple_2[1])


def substract_tuples(
    tuple_1: tuple[int, int], tuple_2: tuple[int, int]
) -> tuple[int, int]:
    return (tuple_1[0] - tuple_2[0], tuple_1[1] - tuple_2[1])


def scalar_times_tuple(scalar: int, in_tuple: tuple[int, int]) -> tuple[int, int]:
    return (scalar * in_tuple[0], scalar * in_tuple[1])


class JungleFace:
    def __init__(self, top_right: tuple[int, int], walls: set[tuple[int, int]]):
        self.walls = walls
        self.top_right = top_right

    def in_map_position(self, location: tuple[int, int]):
        return sum_tuples(location, self.top_right)

    def is_wall(self, location: tuple[int, int]) -> bool:
        return location in self.walls


class JungleCube:
    def __init__(self, terrain: list[str], cube_configuration: bool = False):
        self.faces = []
        self.cube_corners = self.get_top_left(terrain)
        assert len(self.cube_corners) == 6, "Input does not seem to be a cube"
        i = 0
        for corner in self.cube_corners:
            self.faces.append(JungleFace(corner, self.get_walls(terrain, corner)))
            self.faces[i].label = i+1
            i = i+1
        self.set_cube_configuration(cube_configuration)

    @staticmethod
    def compute_cube_size(terrain: list[str]) -> int:
        positions_count = 0
        for line in terrain:
            positions_count += len(line.strip())
        return int((positions_count // 6) ** (1 / 2))

    def get_top_left(self, terrain: list[str]) -> list[tuple[int, int]]:
        self.cube_size = self.compute_cube_size(terrain)
        square_corners = []
        for vertical_index in range(0, len(terrain), self.cube_size):
            for horizontal_index in range(
                0, len(terrain[vertical_index]), self.cube_size
            ):
                if terrain[vertical_index][horizontal_index] != " ":
                    square_corners.append((horizontal_index, vertical_index))
        return square_corners

    def get_walls(
        self, terrain: list[str], top_left: tuple[int, int]
    ) -> set[tuple[int, int]]:
        walls = set()
        for vertical_index in range(top_left[0], top_left[0] + self.cube_size):
            for horizontal_index in range(top_left[1], top_left[1] + self.cube_size):
                if terrain[horizontal_index][vertical_index] == "#":
                    walls.add(
                        substract_tuples((vertical_index, horizontal_index), top_left)
                    )
        return walls

    def set_cube_configuration(self, cube_mode: bool):
        self.restart_configuration()
        for index, corner in enumerate(self.cube_corners):
            for direction_index, direction in enumerate(DIRECTION_VALUES):
                current_tuple = sum_tuples(
                    corner, scalar_times_tuple(self.cube_size, direction)
                )
                if current_tuple in self.cube_corners:
                    self.faces[index].neighbours[
                        direction_index
                    ] = self.cube_corners.index(current_tuple)
        if not cube_mode:
            self.set_cube_wrapping()
        self.fold_cube()

    def set_cube_wrapping(self):
        for index in range(len(self.cube_corners)):
            for direction_index in range(4):
                if self.faces[index].neighbours[direction_index] != -1:
                    continue
                opposite_direction = (direction_index + 2) % 4
                current_index = index
                next_index = self.faces[index].neighbours[opposite_direction]
                while next_index not in [-1, index]:
                    current_index = next_index
                    next_index = self.faces[current_index].neighbours[
                        opposite_direction
                    ]
                self.faces[index].neighbours[direction_index] = current_index

    def fold_cube(self):
        set_up = [False] * len(self.cube_corners)
        while not all(set_up):
            indices = [
                face_index
                for face_index in range(len(set_up))
                if not set_up[face_index]
            ]
            for face_index in indices:
                neighbours = self.faces[face_index].neighbours
                if all([neighbour != -1 for neighbour in neighbours]):
                    set_up[face_index] = True
                    continue
                if neighbours[0] == -1:
                    self.try_fill_right(face_index)
                if neighbours[1] == -1:
                    self.try_fill_bottom(face_index)
                if neighbours[2] == -1:
                    self.try_fill_left(face_index)
                if neighbours[3] == -1:
                    self.try_fill_top(face_index)

    # We fold the cube by the 90 degree angles. We do so until every neighbour is filled.
    # It's just a matter of being careful with the relative angles between the neighbours.
    def try_fill(
        self, face_index: int, direction_1: int, direction_2: int, angle_offset: int
    ) -> bool:
        neighbours = self.faces[face_index].neighbours
        angles = self.faces[face_index].neighbour_angles
        if neighbours[direction_1] == -1:
            return False
        neighbour_1 = neighbours[direction_1]
        neighbour_1_angle = angles[direction_1]
        neighbour_2_index = (direction_2 - neighbour_1_angle) % 4
        neighbour_2 = self.faces[neighbour_1].neighbours[neighbour_2_index]
        if neighbour_2 == -1:
            return False
        neighbour_2_angle = self.faces[neighbour_1].neighbour_angles[neighbour_2_index]
        self.faces[face_index].neighbours[direction_2] = neighbour_2
        self.faces[face_index].neighbour_angles[direction_2] = (
            angle_offset + neighbour_1_angle + neighbour_2_angle
        ) % 4
        return True

    def try_fill_right(self, face_index: int):
        # Fill right by going bottom then right
        if not self.try_fill(face_index, 3, 0, 1):
            # Fill right by going top then right
            self.try_fill(face_index, 1, 0, 3)
        # Similar for the rest of the four directions

    def try_fill_bottom(self, face_index: int):
        if not self.try_fill(face_index, 0, 1, 1):
            self.try_fill(face_index, 2, 1, 3)

    def try_fill_left(self, face_index: int):
        if not self.try_fill(face_index, 3, 2, 3):
            self.try_fill(face_index, 1, 2, 1)

    def try_fill_top(self, face_index: int):
        if not self.try_fill(face_index, 2, 3, 1):
            self.try_fill(face_index, 0, 3, 3)

    def restart_configuration(self):
        self.current_face = 0
        self.position = (0, 0)
        self.direction_index = 0
        for face in self.faces:
            face.neighbours = [-1, -1, -1, -1]
            face.neighbour_angles = [0, 0, 0, 0]

    def wrap_location(
        self, current_location: tuple[int, int], angle: int
    ) -> tuple[int, int]:
        x, y, N = current_location[0], current_location[1], self.cube_size - 1
        match (self.direction_index, angle):
            case 0, 0:
                return (0, y)
            case 0, 1:
                return (y, N)
            case 0, 2:
                return (N, N - y)
            case 0, 3:
                return (N - y, 0)
            case 1, 0:
                return (x, 0)
            case 1, 1:
                return (0, N - x)
            case 1, 2:
                return (N - x, N)
            case 1, 3:
                return (N, x)
            case 2, 0:
                return (N, y)
            case 2, 1:
                return (y, 0)
            case 2, 2:
                return (0, N - y)
            case 2, 3:
                return (N - y, N)
            case 3, 0:
                return (x, N)
            case 3, 1:
                return (N, N - x)
            case 3, 2:
                return (N - x, 0)
            case 3, 3:
                return (0, x)
            case other:
                raise AssertionError

    def is_out_of_bounds(self, position: tuple[int, int]) -> bool:
        return any([element < 0 or element >= self.cube_size for element in position])

    def move(self, number_moves: int):
        direction = DIRECTION_VALUES[self.direction_index]
        for _ in range(number_moves):
            next_position = sum_tuples(self.position, direction)
            if self.is_out_of_bounds(next_position):
                angle = self.faces[self.current_face].neighbour_angles[
                    self.direction_index
                ]
                next_position = self.wrap_location(
                    self.position,
                    angle,
                )
                next_face = self.faces[self.current_face].neighbours[
                    self.direction_index
                ]
                if self.faces[next_face].is_wall(next_position):
                    return None
                self.position = next_position
                self.current_face = next_face
                self.direction_index -= angle
                self.direction_index %= 4
                direction = DIRECTION_VALUES[self.direction_index]
            else:
                if self.faces[self.current_face].is_wall(next_position):
                    return None
                self.position = next_position

    def in_map_location(self) -> tuple[int, int]:
        return self.faces[self.current_face].in_map_position(self.position)

    def password(self) -> int:
        position = self.in_map_location()
        return 1000 * (position[1] + 1) + 4 * (position[0] + 1) + self.direction_index

    def path_move(self, instructions):
        self.direction_index = 0
        direction_changes = {"R": 1, "L": -1}
        for instruction in instructions:
            if isinstance(instruction, int):
                self.move(instruction)
                instlabel = instruction
            else:
                self.direction_index += direction_changes[instruction]
                self.direction_index %= 4
                instlabel = " " + instruction
            print(f"{instlabel}: [{self.faces[self.current_face].label}] {self.position} {self.in_map_location()} {self.direction_index}")


if __name__ == "__main__":
    try:
        file_name = sys.argv[1]
    except IndexError:
        file_name = "input.txt"
    input_data = open(file_name).read().splitlines()
    terrain, moves = input_data[:-2], input_data[-1]

    jungle = JungleCube(terrain)
    instructions = parse_instructions(moves)

    #jungle.path_move(instructions)
    #print(f"The password is {jungle.password()}.")

    print("Running Part 2\ndeividragon solution in python\nCube faces 50")
    jungle.set_cube_configuration(True)
    jungle.path_move(instructions)
    print(f"The password when seen as cube is {jungle.password()}")
